package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/davecgh/go-spew/spew"
	"log"
	"os"
	"time"

	"github.com/golang-coach/mono-repo/constants"
	"github.com/golang-coach/mono-repo/services"
	"github.com/google/go-github/github"
	"github.com/json-iterator/go"
	"golang.org/x/oauth2"

	"github.com/golang-coach/mono-repo/aws"
	"github.com/golang-coach/mono-repo/models"
)

type GithubResponse struct {
	Repository *models.Repository
	err        error
}

func HandleRequest(context context.Context) error {
	db := aws.DBClient(true)
	defer db.Close()

	var repos []models.Repository

	err := db.Find(&repos).Error

	if err != nil {
		return err
	}

	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv(constants.GithubToken)},
	)
	tokenClient := oauth2.NewClient(context, tokenService)
	client := *github.NewClient(tokenClient)
	githubService := services.NewGithub(tokenClient, &client, client.Repositories, context)

	responses := getRepositories(repos, githubService)

	sess := aws.Session()
	queue := aws.NewQueue(sqs.New(sess))
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	spew.Dump(responses)
	for _, res := range responses {
		if res.err != nil {
			fmt.Println(err)
			continue
		}
		if res.Repository != nil {
			message, _ := json.MarshalToString(res.Repository)
			if err := queue.Send(message); err != nil {
				fmt.Print(err)
			}
		}
	}
	return nil
}

func getRepositories(repos []models.Repository, githubService services.Github) []*GithubResponse {
	ch := make(chan *GithubResponse)
	var responses []*GithubResponse
	for _, repo := range repos {
		go func(repository models.Repository) {
			fmt.Print(repository.Name)
			lastCommitInfo, err := githubService.GetLastCommitInfo(repo.Owner, repo.Name)
			spew.Dump(lastCommitInfo)
			if err != nil {
				ch <- &GithubResponse{
					err: err,
				}
			} else {
				lastCommitDate := lastCommitInfo.Commit.Committer.GetDate()
				if lastCommitDate.After(repo.ProcessedAt) {
					ch <- &GithubResponse{
						Repository: &repository,
					}
				} else {
					ch <- &GithubResponse{}
				}
			}
		}(repo)
	}
	for {
		select {
		case res := <-ch:
			responses = append(responses, res)
			if len(responses) == len(repos) {
				return responses
			}
		case <-time.After(10 * time.Minute):
			log.Fatalln("Timeout")
			return responses
		}
	}
	return responses
}

func main() {
	lambda.Start(HandleRequest)
}

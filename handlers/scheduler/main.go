package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang-coach/mono-repo/aws"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/json-iterator/go"
	"golang.org/x/oauth2"

	"github.com/golang-coach/mono-repo/constants"
	"github.com/golang-coach/mono-repo/models"
	"github.com/golang-coach/mono-repo/services"
)

type GithubResponse struct {
	Repository *models.Repository
	err        error
}

const LIMIT = 20

func HandleRequest(context context.Context) error {
	db := aws.DBClient(false)
	defer db.Close()

	var repos []models.Repository

	yesterday := time.Now().Add(-24 * time.Hour)
	err := db.
		Where("updated_at < ?", yesterday).
		Or("updated_at is NULL").
		Limit(LIMIT).
		Find(&repos).
		Error

	if err != nil {
		return err
	}

	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv(constants.GithubToken)},
	)
	tokenClient := oauth2.NewClient(context, tokenService)
	client := *github.NewClient(tokenClient)
	githubService := services.NewGithub(tokenClient, &client, client.Repositories, context)

	responses := getRepositories(db, repos, githubService)

	sess := aws.Session()
	queue := aws.NewQueue(sqs.New(sess))
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	for _, res := range responses {
		if res.err != nil {
			fmt.Println("Error " + res.err.Error())
			continue
		}
		if res.Repository != nil {
			fmt.Println("Repo " + res.Repository.Name + " Full Name " + res.Repository.FullName)
			message, _ := json.MarshalToString(res.Repository)
			if err := queue.Send(message); err != nil {
				fmt.Println("Error " + err.Error())
			}
		}
	}
	return nil
}

func getRepositories(db *gorm.DB, repos []models.Repository, githubService services.Github) []*GithubResponse {
	var responses []*GithubResponse
	now := time.Now()
	for index, repository := range repos {
		repository.UpdatedAt = now
		err := db.
			Save(&repository).
			Error
		if err != nil {
			responses = append(responses, &GithubResponse{
				err: err,
			})
			continue
		}
		responses = append(responses, &GithubResponse{Repository: nil})
		lastCommitInfo, err := githubService.GetLastCommitInfo(repository.Owner, repository.Name)
		if err != nil {
			responses = append(responses, &GithubResponse{
				err: err,
			})
		} else {
			lastCommitDate := lastCommitInfo.Commit.Committer.GetDate()
			if lastCommitDate.After(repository.ProcessedAt) {
				responses = append(responses, &GithubResponse{
					Repository: &repos[index],
				})
			} else {
				responses = append(responses, &GithubResponse{})
			}
		}
	}
	return responses
}

func main() {
	lambda.Start(HandleRequest)
	//HandleRequest(context.Background())
}

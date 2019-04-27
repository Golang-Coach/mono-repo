package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-coach/mono-repo/aws"
	"github.com/golang-coach/mono-repo/constants"
	"github.com/golang-coach/mono-repo/models"
	"github.com/golang-coach/mono-repo/services"
	"github.com/google/go-github/github"
	"github.com/json-iterator/go"
	"golang.org/x/oauth2"
)

func HandleRequest(context context.Context, request events.SQSEvent) error {
	if len(request.Records) > 0 {
		repository := models.Repository{}
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		err := json.UnmarshalFromString(request.Records[0].Body, &repository)
		if err != nil {
			return err
		}

		tokenService := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv(constants.GithubToken)},
		)
		tokenClient := oauth2.NewClient(context, tokenService)
		client := *github.NewClient(tokenClient)
		githubService := services.NewGithub(tokenClient, &client, client.Repositories, context)

		repo, err := githubService.GetUpdatedRepositoryInfo(repository)
		if err != nil {
			return err
		}

		db := aws.DBClient(true)
		defer db.Close()

		return db.Save(&repo).Error
	}
	err := errors.New("no message found")
	return err
}

func main() {
	lambda.Start(HandleRequest)
}

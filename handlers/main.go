package main

import (
	"context"
	"fmt"
	"github.com/golang-coach/mono-repo/models"
	"github.com/golang-coach/mono-repo/services"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
)

func main() {
	repository := models.Repository{
		Name:  "env",
		Owner: "codingconcepts",
	}
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//err := json.UnmarshalFromString(request.Records[0].Body, &repository)
	//if err != nil {
	//	return err
	//}

	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("7047fe0802e602011804f0c15779b48dcb5388df")},
	)
	context1 := context.Background()
	tokenClient := oauth2.NewClient(context1, tokenService)
	client := *github.NewClient(tokenClient)
	githubService := services.NewGithub(tokenClient, &client, client.Repositories, context1)

	fmt.Println("Processing repo " + repository.Name)
	repo, err := githubService.GetUpdatedRepositoryInfo(repository)
	if err != nil {
		panic(err)
	}
	fmt.Println(repo.Name)
}

package services

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	githubReadMe "github.com/Golang-Coach/github-readme"
	"github.com/golang-coach/mono-repo/interfaces"
	"github.com/golang-coach/mono-repo/models"
	"github.com/google/go-github/github"
)

type IClient interface {
	RateLimits(ctx context.Context) (*github.RateLimits, *github.Response, error)
}

type Github struct {
	httpClient         *http.Client
	client             IClient
	repositoryServices interfaces.IRepositoryServices
	context            context.Context
}

func NewGithub(httpClient *http.Client, client IClient, repositoryServices interfaces.IRepositoryServices, context context.Context) Github {
	return Github{
		httpClient:         httpClient,
		client:             client,
		repositoryServices: repositoryServices,
		context:            context,
	}
}

func (service Github) GetRepositoryInfo(owner string, repositoryName string) (*models.Repository, error) {
	repo, _, err := service.repositoryServices.Get(service.context, owner, repositoryName)
	if err != nil {
		return nil, err
	}
	repositoryInfo := &models.Repository{
		Name:        *repo.Name,
		Owner:       strings.Split(*repo.FullName, "/")[0],
		FullName:    *repo.FullName,
		Description: *repo.Description,
		Forks:       *repo.ForksCount,
		Stars:       *repo.StargazersCount,
	}
	return repositoryInfo, nil
}

func (service Github) GetLastCommitInfo(owner string, repositoryName string) (*github.RepositoryCommit, error) {
	commitInfo, _, err := service.repositoryServices.ListCommits(service.context, owner, repositoryName, nil)
	if err != nil {
		return nil, err
	}
	return commitInfo[0], nil
}

func (service Github) GetReadMe(owner string, repositoryName string) (string, error) {
	readme, _, err := service.repositoryServices.GetReadme(service.context, owner, repositoryName, nil)
	if err != nil {
		return "", err
	}

	// get content
	return readme.GetContent()
}

func (service Github) GetRateLimitInfo() (*github.RateLimits, error) {
	rateLimitInfo, _, err := service.client.RateLimits(service.context)
	return rateLimitInfo, err
}

func (service Github) GetUpdatedRepositoryInfo(repositoryInfo models.Repository) (*models.Repository, error) {
	// Call last update information from Github API
	if len(strings.TrimSpace(repositoryInfo.Owner)) == 0 || len(strings.TrimSpace(repositoryInfo.Name)) == 0 {
		return nil, errors.New("Repository Name is incorrect")
	}

	// check with existing caller api
	lastCommitInfo, err := service.GetLastCommitInfo(repositoryInfo.Owner, repositoryInfo.Name)
	if err != nil {
		return nil, err
	}

	if lastCommitInfo.Commit.Committer.GetDate().Equal(repositoryInfo.UpdatedAt) {
		return nil, nil
	}

	newRepositoryInfo, err := service.GetRepositoryInfo(repositoryInfo.Owner, repositoryInfo.Name)
	if err != nil {
		return nil, err
	}

	newRepositoryInfo.ID = repositoryInfo.ID
	newRepositoryInfo.UpdatedAt = lastCommitInfo.Commit.Committer.GetDate()
	avatarUrl := ""
	committerName := *lastCommitInfo.Commit.Committer.Name
	if lastCommitInfo.Author != nil && lastCommitInfo.Author.AvatarURL != nil {
		avatarUrl = *lastCommitInfo.Author.AvatarURL
		committerName = *lastCommitInfo.Commit.Author.Name
	}
	newRepositoryInfo.User = models.User{
		Name:       committerName,
		UserName:   committerName,
		ProfileUrl: avatarUrl,
	}

	newRepositoryInfo.LastUpdatedBy = committerName

	githubReadMeClient := githubReadMe.NewGithub(service.httpClient)
	content, err := githubReadMeClient.GetReadme(service.context, repositoryInfo.Owner, repositoryInfo.Name)

	if err != nil {
		return nil, err
	}
	newRepositoryInfo.ReadMe = content
	newRepositoryInfo.Processed = true
	newRepositoryInfo.ProcessedAt = time.Now()

	// data is same ignore, else update data
	return newRepositoryInfo, nil
}

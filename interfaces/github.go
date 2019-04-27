package interfaces

import (
	"context"

	"github.com/golang-coach/mono-repo/models"
	"github.com/google/go-github/github"
)

type IRepositoryServices interface {
	Get(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error)
	ListCommits(ctx context.Context, owner, repo string, opt *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error)
	GetReadme(ctx context.Context, owner, repo string, opt *github.RepositoryContentGetOptions) (*github.RepositoryContent, *github.Response, error)
}

type IGithub interface {
	GetRepositoryInfo(owner string, repositoryName string) (*models.Repository, error)
	GetLastCommitInfo(owner string, repositoryName string) (*github.RepositoryCommit, error)
	GetReadMe(owner string, repositoryName string) (string, error)
	GetRateLimitInfo() (*github.RateLimits, error)
	GetUpdatedRepositoryInfo(repositoryInfo models.Repository) (*models.Repository, error)
}

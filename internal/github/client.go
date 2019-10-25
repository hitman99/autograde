package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Client interface {
	CheckFork(ctx context.Context, owner, repo string) (bool, error)
	CheckBuildAction() (bool, error)
}

type ghClient struct {
	client *githubv4.Client
}

func MustNewClient(token string) Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	return &ghClient{
		client: client,
	}
}

func (g *ghClient) CheckFork(ctx context.Context, owner, repo string) (bool, error) {
	var fork struct {
		Repository struct {
			IsFork bool
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	err := g.client.Query(ctx, &fork, map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(repo),
	})
	if err != nil {
		return false, err
	}
	return fork.Repository.IsFork, nil
}

func (g *ghClient) CheckBuildAction() (bool, error) {
	return false, nil
}

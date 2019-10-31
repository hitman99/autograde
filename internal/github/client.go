package github

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Client interface {
	CheckFork(ctx context.Context, owner, repo string) (bool, error)
	CheckBuildAction(ctx context.Context, owner, repo, actionName string) (bool, error)
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
		if err.Error() == fmt.Sprintf("Could not resolve to a Repository with the name '%s'.", repo) {
			return false, nil
		} else {
			return false, err
		}
	}
	return fork.Repository.IsFork, nil
}

func (g *ghClient) CheckBuildAction(ctx context.Context, owner, repo, actionName string) (bool, error) {
	var query struct {
		Repository struct {
			Object struct {
				Blob struct {
					Id githubv4.ID
				} `graphql:"... on Blob"`
			} `graphql:"object(expression: $expr)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	err := g.client.Query(ctx, &query, map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(repo),
		"expr":  githubv4.String("master:.github/workflows/release.yml"),
	})
	if err != nil {
		if err.Error() == fmt.Sprintf("Could not resolve to a Repository with the name '%s'.", repo) {
			return false, nil
		} else {
			return false, err
		}
	}
	return query.Repository.Object.Blob.Id != nil, nil
}

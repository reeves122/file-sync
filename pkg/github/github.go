package github

import (
	"context"
	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
	"strings"
)

func GetClient(token string) *github.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	return github.NewClient(httpClient)
}

func CreatePullRequest(client *github.Client, owner, repo, title, head, base string) error {
	_, _, err := client.PullRequests.Create(context.Background(), owner, repo, &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(head),
		Base:  github.String(base),
	})
	if err != nil {
		if strings.Contains(err.Error(), "A pull request already exists") {
			return nil
		}
		return err
	}
	return nil
}

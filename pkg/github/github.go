package github

import (
	"context"
	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

func GetClient(token string) *github.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	return github.NewClient(httpClient)
}

func CreatePullRequest(client *github.Client) error {
	_, _, err := client.PullRequests.Create(context.Background(), "reeves122", "file-sync", &github.NewPullRequest{
		Title: github.String("file sync"),
		Head:  github.String("file-sync"),
		Base:  github.String("main"),
	})
	if err != nil {
		return err
	}
	return nil
}

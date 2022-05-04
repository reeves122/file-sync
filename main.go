package main

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"time"
)

const sourceRepo = "https://github.com/champ-oss/terraform-module-template"
const sourceRepoTmp = "/tmp/source"

func main() {
	_, err := git.PlainClone(sourceRepoTmp, false, &git.CloneOptions{
		URL:      sourceRepo,
		Progress: os.Stdout,
	})
	if err != nil {
		panic(err)
	}

	r, err := git.PlainOpen("./")
	if err != nil {
		panic(err)
	}

	w, err := r.Worktree()
	if err != nil {
		panic(err)
	}

	fmt.Println("checkout test branch")
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("test"),
		Create: true,
		Keep:   true,
	})
	if err != nil {
		panic(err)
	}

	input, err := ioutil.ReadFile("/tmp/source/.tflint.hcl")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(".tflint.hcl", input, 0644)
	if err != nil {
		panic(err)
	}

	_, err = w.Add(".tflint.hcl")
	if err != nil {
		panic(err)
	}

	status, err := w.Status()
	if err != nil {
		panic(err)
	}

	fmt.Printf("git status: \n%s\n", status)

	fmt.Println("making commit")
	commit, err := w.Commit("test commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "file-sync",
			Email: "file-sync@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("getting commit")
	obj, err := r.CommitObject(commit)
	if err != nil {
		panic(err)
	}
	fmt.Println(obj)

	ref, err := r.Head()
	fmt.Println("ref:", ref)
	if err != nil {
		panic(err)
	}

	fmt.Println("pushing")
	err = r.Push(&git.PushOptions{
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: "testuser",
			Password: os.Getenv("GITHUB_TOKEN"),
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("creating github client")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)
	_, resp, err := client.PullRequests.Create(context.Background(), "reeves122", "file-sync", &github.NewPullRequest{
		Title: github.String("file sync"),
		Head:  github.String("test"),
		Base:  github.String("main"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

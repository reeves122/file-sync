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

func main() {
	//r, err := git.PlainClone("./foo", false, &git.CloneOptions{
	//	URL:      "https://github.com/go-git/go-git",
	//	Progress: os.Stdout,
	//})
	//r, err := git.PlainClone("./foo", false, &git.CloneOptions{
	//	Auth: &http.BasicAuth{
	//		Username: "abc123",
	//		Password: os.Getenv("GITHUB_TOKEN"),
	//	},
	//	URL:      "./",
	//	Progress: os.Stdout,
	//})
	r, err := git.PlainOpen("./")
	if err != nil {
		panic(err)
	}

	//err = r.DeleteRemote("origin")
	//if err != nil {
	//	panic(err)
	//}
	//
	//_, err = r.CreateRemote(&config.RemoteConfig{
	//	Name: "origin",
	//	URLs: []string{fmt.Sprintf("https://filesync:%s@github.com/reeves122/file-sync", os.Getenv("GITHUB_TOKEN"))},
	//})
	//if err != nil {
	//	panic(err)
	//}

	remotes, err := r.Remotes()
	if err != nil {
		panic(err)
	}
	fmt.Println("remotes:")
	for _, remote := range remotes {
		fmt.Println(remote)
	}

	//err = r.CreateBranch(&config.Branch{
	//	Name: "test",
	//})
	//if err != nil {
	//	panic(err)
	//}

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

	err = ioutil.WriteFile("./test-file.txt", []byte("test"), 0644)
	if err != nil {
		panic(err)
	}

	_, err = w.Add("test-file.txt")

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
		//RemoteName: "origin",
		//RefSpecs: []config.RefSpec{"refs/heads/test:refs/heads/test"},
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

	fmt.Println("listing branches")
	client := github.NewClient(tc)
	branches, resp, err := client.Repositories.ListBranches(context.Background(), "reeves122", "file-sync", &github.BranchListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
	fmt.Println(branches)
}

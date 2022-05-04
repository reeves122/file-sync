package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"io/ioutil"
	"time"
)

func main() {
	//r, err := git.PlainClone("./foo", false, &git.CloneOptions{
	//	URL:      "https://github.com/go-git/go-git",
	//	Progress: os.Stdout,
	//})
	r, err := git.PlainOpen("./")
	if err != nil {
		panic(err)
	}

	err = r.CreateBranch(&config.Branch{
		Name: "test",
	})
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("test-file.txt", []byte("test"), 0644)
	if err != nil {
		panic(err)
	}

	w, err := r.Worktree()
	if err != nil {
		panic(err)
	}

	_, err = w.Add("test-file.txt")

	status, err := w.Status()
	if err != nil {
		panic(err)
	}

	fmt.Printf("git status: \n%s\n", status)

	commit, err := w.Commit("test commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "file-sync",
			Email: "file-sync@example.com",
			When:  time.Now(),
		},
	})

	obj, err := r.CommitObject(commit)
	if err != nil {
		panic(err)
	}

	fmt.Println(obj)

	err = r.Push(&git.PushOptions{})
	if err != nil {
		panic(err)
	}
}

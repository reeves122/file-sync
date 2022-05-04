package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

func main() {
	//r, err := git.PlainClone("./foo", false, &git.CloneOptions{
	//	URL:      "https://github.com/go-git/go-git",
	//	Progress: os.Stdout,
	//})
	r, err := git.PlainOpen("./")
	if err != nil {
		fmt.Println(err)
	}

	err = r.CreateBranch(&config.Branch{
		Name: "test",
	})

	if err != nil {
		fmt.Println(err)
	}

	w, err := r.Worktree()
	status, err := w.Status()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("git status: \n%s\n", status)
}

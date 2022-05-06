package main

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v44/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const localRepoDir = "./"
const branchName = "file-sync"
const commitMsg = "file-sync"
const user = "file-sync"
const email = "file-sync@example.com"

var files = []string{
	".tflint.hcl",
	".github/CODEOWNERS",
	".github/workflows/release.yml",
}

func main() {
	log.SetLevel(log.DebugLevel)

	sourceDir, err := cloneSourceRepo("https://github.com/champ-oss/terraform-module-template")
	if err != nil {
		log.Fatal(err)
	}

	repo, err := openLocalRepo(localRepoDir)
	if err != nil {
		log.Fatal(err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		log.Fatal("error getting the worktree for the local repository", err)
	}

	err = repo.Fetch(&git.FetchOptions{
		//RemoteName:      "",
		//RefSpecs:        nil,
		//Depth:           0,
		Progress: os.Stdout,
		//Tags:            0,
	})
	if err != nil {
		log.Fatal("error getting the worktree for the local repository", err)
	}

	if err := checkOutBranch(worktree, branchName); err != nil {
		log.Fatal(err)
	}

	if err := copySourceFiles(files, sourceDir, localRepoDir); err != nil {
		log.Fatal(err)
	}

	modified, err := isWorktreeModified(worktree)
	if err != nil {
		log.Fatal(err)
	}
	if !modified {
		log.Info("all files are up to date")
		os.Exit(0)
	}

	if err := gitAddFiles(files, worktree); err != nil {
		log.Fatal(err)
	}

	hash, err := createCommit(worktree, commitMsg, user, email)
	if err != nil {
		log.Fatal(err)
	}

	commit, err := repo.CommitObject(hash)
	fmt.Println(commit.Files())
	fmt.Println(commit.String())

	if err := gitPush(repo, user, os.Getenv("GITHUB_TOKEN")); err != nil {
		log.Fatal(err)
	}

	fmt.Println("creating github client")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)
	_, resp, err := client.PullRequests.Create(context.Background(), "reeves122", "file-sync", &github.NewPullRequest{
		Title: github.String("file sync"),
		Head:  github.String("file-sync"),
		Base:  github.String("main"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func cloneSourceRepo(sourceRepo string) (dir string, err error) {
	log.Debug("Creating temp directory for source repository")
	dir, _ = ioutil.TempDir("", "source")

	log.Infof("Cloning source repository %s to %s", sourceRepo, dir)
	if _, err := git.PlainClone(dir, false, &git.CloneOptions{URL: sourceRepo, Progress: os.Stdout}); err != nil {
		log.Error("error cloning source repository")
		return dir, err
	}
	return dir, nil
}

func openLocalRepo(path string) (repo *git.Repository, err error) {
	log.Infof("Opening local repository: %s", path)
	repo, err = git.PlainOpen(path)
	if err != nil {
		log.Error("error opening local directory as git repository")
		return nil, err
	}
	return repo, nil
}

func copySourceFiles(files []string, sourceDir, destDir string) error {
	for _, f := range files {
		sourcePath := filepath.Join(sourceDir, f)
		destPath := filepath.Join(destDir, f)
		log.Debugf("Copying %s to %s", sourcePath, destPath)
		if err := copyFile(sourcePath, destPath); err != nil {
			log.Error("error copying files from source")
			return err
		}
	}
	return nil
}

func copyFile(source, dest string) error {
	input, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	if baseDir, _ := filepath.Split(dest); baseDir != "" {
		if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(dest, input, 0644)
	return err
}

func isWorktreeModified(worktree *git.Worktree) (bool, error) {
	status, err := worktree.Status()
	if err != nil {
		log.Error("error running git status")
		return false, err
	}
	log.Info("git status")
	fmt.Print(status)
	return !status.IsClean(), nil
}

func checkOutBranch(worktree *git.Worktree, branchName string) error {
	log.Infof("Checking out branch: %s", branchName)
	err := worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
		Create: true,
		Keep:   true,
	})

	if err != nil && strings.Contains(err.Error(), "already exists") {

		log.Debugf("Branch already exists. Switching to existing branch: %s", branchName)
		if err := worktree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(branchName),
			Create: false,
			Keep:   true,
		}); err != nil {
			return err
		} else {
			return nil
		}
	}

	if err != nil {
		log.Errorf("Error checking out branch: %s", branchName)
		return err
	}
	return nil
}

func gitAddFiles(files []string, worktree *git.Worktree) error {
	for _, f := range files {
		log.Debugf("git add %s", f)
		if _, err := worktree.Add(f); err != nil {
			log.Errorf("error running git add %s", f)
			return err
		}
	}
	return nil
}

func createCommit(worktree *git.Worktree, msg, name, email string) (plumbing.Hash, error) {
	log.Infof("Creating commit as user %s", name)
	hash, err := worktree.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
	})
	if err != nil {
		log.Error("error creating commit")
		return hash, err
	}
	log.Infof("Created commit: %s", hash)
	return hash, nil
}

func gitPush(repo *git.Repository, username, password string) error {
	log.Info("Running git push")
	err := repo.Push(&git.PushOptions{
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		log.Error("error running git push")
		return err
	}
	log.Infof("Successfully pushed")
	return nil
}

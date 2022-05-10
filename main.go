package main

import (
	"github.com/champ-oss/file-sync/pkg/common"
	"github.com/champ-oss/file-sync/pkg/git/cli"
	"github.com/champ-oss/file-sync/pkg/github"
	log "github.com/sirupsen/logrus"
	"os"
)

const localRepoDir = "./"
const branchName = "file-sync"
const commitMsg = "file-sync"
const user = "file-sync"
const email = "file-sync@example.com"

var files = []string{
	".tflint.hcl",
	"test/src/go.mod",
	"examples/complete/main.tf",
	".github/CODEOWNERS",
	".github/workflows/release.yml",
}

func main() {
	log.SetLevel(log.DebugLevel)

	owner := common.GetOwner()
	repo := common.GetRepo()
	//repo := os.Getenv("GITHUB_REPOSITORY")
	token := os.Getenv("INPUT_TOKEN")

	_, err := common.RunCommand(localRepoDir, "ls", "-l")
	if err != nil {
		log.Fatal(err)
	}

	sourceDir, err := cli.Clone("champ-oss/terraform-module-template", token)
	if err != nil {
		log.Fatal(err)
	}

	//destDir, err := cli.Clone(repo, token)
	//if err != nil {
	//	log.Fatal(err)
	//}

	err = cli.SetAuthor(localRepoDir, user, email)
	if err != nil {
		panic(err)
	}

	err = cli.Fetch(localRepoDir)
	if err != nil {
		panic(err)
	}

	err = cli.Branch(localRepoDir, branchName)
	if err != nil {
		panic(err)
	}

	err = cli.Checkout(localRepoDir, branchName)
	if err != nil {
		panic(err)
	}

	err = cli.Reset(localRepoDir, branchName)
	if err != nil {
		panic(err)
	}

	if err := common.CopySourceFiles(files, sourceDir, localRepoDir); err != nil {
		log.Fatal(err)
	}

	if modified := cli.AnyModified(localRepoDir, files); modified == false {
		log.Info("all files are up to date")
		os.Exit(0)
	}

	for _, f := range files {
		err = cli.Add(localRepoDir, f)
		if err != nil {
			panic(err)
		}
	}

	err = cli.Commit(localRepoDir, commitMsg)
	if err != nil {
		log.Fatal(err)
	}

	_, err = common.RunCommand(localRepoDir, "git", "remote", "-v")
	if err != nil {
		log.Fatal(err)
	}

	err = cli.Push(localRepoDir, branchName)
	if err != nil {
		log.Fatal(err)
	}

	client := github.GetClient(token)
	err = github.CreatePullRequest(client, owner, repo, "file-sync", branchName, "main")
	if err != nil {
		log.Fatal(err)
	}
}

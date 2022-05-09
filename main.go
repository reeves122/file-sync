package main

import (
	"github.com/champ-oss/file-sync/pkg/common"
	"github.com/champ-oss/file-sync/pkg/git/cli"
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
	sourceDir, err := cli.Clone("https://github.com/champ-oss/terraform-module-template")
	if err != nil {
		log.Fatal(err)
	}

	err = cli.SetAuthor(localRepoDir, user, email)
	if err != nil {
		panic(err)
	}
	//
	//_, err = common.RunCommand(localRepoDir, "git", "remote", "remove", "origin")
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = common.RunCommandNoLog(localRepoDir, "git", "remote", "add", "origin", fmt.Sprintf("https://%s@github.com/reeves122/file-sync.git", os.Getenv("FILE_SYNC_PAT")))
	//if err != nil {
	//	panic(err)
	//}

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

	//_, err = common.RunCommand(localRepoDir, "git", "push", fmt.Sprintf("https://%s@github.com/reeves122/file-sync.git", os.Getenv("FILE_SYNC_PAT")))
	//if err != nil {
	//	log.Fatal(err)
	//}
	err = cli.Push(localRepoDir, branchName)
	if err != nil {
		log.Fatal(err)
	}

}

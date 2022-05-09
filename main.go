package main

import (
	log "github.com/sirupsen/logrus"
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
	//
	//sourceDir, err := cli.CloneSourceRepo("https://github.com/champ-oss/terraform-module-template")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//if err := checkOutBranch(worktree, branchName); err != nil {
	//	log.Fatal(err)
	//}
	//
	//if err := copySourceFiles(files, sourceDir, localRepoDir); err != nil {
	//	log.Fatal(err)
	//}
	//
	//modified, err := isWorktreeModified(worktree)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if !modified {
	//	log.Info("all files are up to date")
	//	os.Exit(0)
	//}
	//
	//if err := gitAddFiles(files, worktree); err != nil {
	//	log.Fatal(err)
	//}
	//
	//hash, err := createCommit(worktree, commitMsg, user, email)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//commit, err := repo.CommitObject(hash)
	//fmt.Println(commit.Files())
	//fmt.Println(commit.String())
	//
	//if err := gitPush(repo, user, os.Getenv("GITHUB_TOKEN")); err != nil {
	//	log.Fatal(err)
	//}

}

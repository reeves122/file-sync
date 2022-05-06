package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var testBranch = "integration-test"
var testGithubToken = os.Getenv("GITHUB_TOKEN")
var testUser = "integration-test"
var testEmail = "integration-test@example.com"

func init() {
	log.SetLevel(log.DebugLevel)
}

func removeDir(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		log.Error(err)
	}
}

func Test_cloneSourceRepo_Success(t *testing.T) {
	repoDir, err := cloneSourceRepo("https://github.com/git-fixtures/basic.git")
	defer removeDir(repoDir)
	assert.NoError(t, err)

	_, err = os.Stat(repoDir)
	assert.NoError(t, err)
}

func Test_cloneSourceRepo_Error(t *testing.T) {
	repoDir, err := cloneSourceRepo("https://localhost/not-a-repo")
	defer removeDir(repoDir)
	assert.Error(t, err)

	_, err = os.Stat(repoDir)
	assert.NoError(t, err)
}

func Test_openLocalRepo_Success(t *testing.T) {
	repo, err := openLocalRepo("./")
	assert.NoError(t, err)
	assert.NotNil(t, repo)
}

func Test_openLocalRepo_Error(t *testing.T) {
	repo, err := openLocalRepo(os.TempDir())
	assert.Error(t, err)
	assert.Nil(t, repo)
}

func Test_copyFile_Success(t *testing.T) {
	// Write a test file to a test source directory
	sourceDir, _ := ioutil.TempDir("", "source")
	defer removeDir(sourceDir)
	sourceFile := filepath.Join(sourceDir, "test.txt")
	err := ioutil.WriteFile(sourceFile, []byte("test"), 0644)
	assert.NoError(t, err)

	// Copy file and check if it exists in destination
	assert.NoError(t, copyFile(sourceFile, "test.txt"))
	defer os.Remove("test.txt")
	_, err = os.Stat("test.txt")
	assert.NoError(t, err)
}

func Test_copyFile_Bad_Source(t *testing.T) {
	assert.Error(t, copyFile("/foo/invalid.txt", "/foo/invalid.txt"))
}

func Test_copyFile_Bad_Destination(t *testing.T) {
	// Write a test file to a test source directory
	sourceDir, _ := ioutil.TempDir("", "source")
	defer removeDir(sourceDir)
	sourceFile := filepath.Join(sourceDir, "test.txt")
	err := ioutil.WriteFile(sourceFile, []byte("test"), 0644)
	assert.NoError(t, err)
	_, err = os.Stat(sourceFile)
	assert.NoError(t, err)

	// Assert an error when copying to a bad destination
	assert.Error(t, copyFile(sourceFile, "/foo/invalid.txt"))
}

func Test_copyFile_Create_Dir(t *testing.T) {
	// Write a test file to a test source directory
	sourceDir, _ := ioutil.TempDir("", "source")
	defer removeDir(sourceDir)
	sourceFile := filepath.Join(sourceDir, "test.txt")
	err := ioutil.WriteFile(sourceFile, []byte("test"), 0644)
	assert.NoError(t, err)
	_, err = os.Stat(sourceFile)
	assert.NoError(t, err)

	// Create a test destination directory
	destDir, _ := ioutil.TempDir("", "dest")
	defer removeDir(destDir)

	// Use a nested directory that does not exist
	destFile := filepath.Join(destDir, "somedirectory", "test.txt")
	// Copy file and check if it exists in destination
	assert.NoError(t, copyFile(sourceFile, destFile))
	_, err = os.Stat(destFile)
	assert.NoError(t, err)
}

func Test_copySourceFiles_Success(t *testing.T) {
	// Write a test file to a test source directory
	sourceDir, _ := ioutil.TempDir("", "source")
	defer removeDir(sourceDir)
	sourceFile := filepath.Join(sourceDir, "test.txt")
	err := ioutil.WriteFile(sourceFile, []byte("test"), 0644)
	assert.NoError(t, err)
	_, err = os.Stat(sourceFile)
	assert.NoError(t, err)

	// Create a test destination directory
	destDir, _ := ioutil.TempDir("", "dest")
	defer removeDir(destDir)

	assert.NoError(t, copySourceFiles([]string{"test.txt"}, sourceDir, destDir))
}

func Test_copySourceFiles_Error(t *testing.T) {
	assert.Error(t, copySourceFiles([]string{"test.txt"}, "/foo", "/foo"))
}

func Test_isWorktreeModified_Modified(t *testing.T) {
	// Clone and open an example git repository
	repoDir, err := cloneSourceRepo("https://github.com/git-fixtures/basic.git")
	defer removeDir(repoDir)
	repo, err := openLocalRepo(repoDir)
	assert.NoError(t, err)

	// Write modifications to an existing file in the git repository
	err = ioutil.WriteFile(filepath.Join(repoDir, "LICENSE"), []byte("test"), 0644)
	assert.NoError(t, err)

	worktree, err := repo.Worktree()
	assert.NoError(t, err)

	modified, err := isWorktreeModified(worktree)
	assert.NoError(t, err)
	assert.True(t, modified)
}

func Test_isWorktreeModified_New(t *testing.T) {
	// Clone and open an example git repository
	repoDir, err := cloneSourceRepo("https://github.com/git-fixtures/basic.git")
	defer removeDir(repoDir)
	repo, err := openLocalRepo(repoDir)
	assert.NoError(t, err)

	// Write a new file in the git repository
	err = ioutil.WriteFile(filepath.Join(repoDir, "foo-new-file.txt"), []byte("test"), 0644)
	assert.NoError(t, err)

	worktree, err := repo.Worktree()
	assert.NoError(t, err)

	modified, err := isWorktreeModified(worktree)
	assert.NoError(t, err)
	assert.True(t, modified)
}

func Test_isWorktreeModified_Clean(t *testing.T) {
	repoDir, err := cloneSourceRepo("https://github.com/git-fixtures/basic.git")
	defer removeDir(repoDir)
	repo, err := openLocalRepo(repoDir)
	assert.NoError(t, err)

	worktree, err := repo.Worktree()
	assert.NoError(t, err)

	modified, err := isWorktreeModified(worktree)
	assert.NoError(t, err)
	assert.False(t, modified)
}

func Test_checkOutBranch_New(t *testing.T) {
	repoDir, err := cloneSourceRepo("https://github.com/git-fixtures/basic.git")
	defer removeDir(repoDir)
	repo, err := openLocalRepo(repoDir)
	if err != nil {
		log.Fatalf("error opening local directory as git repository \n%s", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		panic(err)
	}

	assert.NoError(t, checkOutBranch(worktree, "foo"))
}

func Test_checkOutBranch_Existing(t *testing.T) {
	repoDir, err := cloneSourceRepo("https://github.com/git-fixtures/basic.git")
	defer removeDir(repoDir)
	repo, err := openLocalRepo(repoDir)
	if err != nil {
		log.Fatalf("error opening local directory as git repository \n%s", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		panic(err)
	}

	assert.NoError(t, checkOutBranch(worktree, "master"))
}

func Test_gitAddFiles_Success(t *testing.T) {
	// Clone and open an example git repository
	repoDir, err := cloneSourceRepo("https://github.com/git-fixtures/basic.git")
	defer removeDir(repoDir)
	repo, err := openLocalRepo(repoDir)
	assert.NoError(t, err)

	// Write a new file in the git repository
	err = ioutil.WriteFile(filepath.Join(repoDir, "foo-new-file.txt"), []byte("test"), 0644)
	assert.NoError(t, err)

	worktree, err := repo.Worktree()
	assert.NoError(t, err)
	assert.NoError(t, gitAddFiles([]string{"foo-new-file.txt"}, worktree))
}

func Test_gitAddFiles_Error(t *testing.T) {
	// Clone and open an example git repository
	repoDir, err := cloneSourceRepo("https://github.com/git-fixtures/basic.git")
	defer removeDir(repoDir)
	repo, err := openLocalRepo(repoDir)
	assert.NoError(t, err)

	worktree, err := repo.Worktree()
	assert.NoError(t, err)
	assert.Error(t, gitAddFiles([]string{"foo-new-file.txt"}, worktree))
}

func Test_createCommit_Success(t *testing.T) {
	// Clone and open an example git repository
	repoDir, err := cloneSourceRepo("https://github.com/git-fixtures/basic.git")
	defer removeDir(repoDir)
	repo, err := openLocalRepo(repoDir)
	assert.NoError(t, err)

	// Write a new file in the git repository
	err = ioutil.WriteFile(filepath.Join(repoDir, "foo-new-file.txt"), []byte("test"), 0644)
	assert.NoError(t, err)

	// Git add and commit
	worktree, err := repo.Worktree()
	assert.NoError(t, err)
	assert.NoError(t, gitAddFiles([]string{"foo-new-file.txt"}, worktree))

	hash, err := createCommit(worktree, "test commit", "example-user", "example-user@example.com")
	assert.NoError(t, err)
	assert.Len(t, hash.String(), 40)

	// Validate details on the commit
	commit, err := repo.CommitObject(hash)
	assert.Equal(t, "test commit", commit.Message)
	assert.Equal(t, "example-user", commit.Author.Name)
	assert.Equal(t, "example-user@example.com", commit.Author.Email)
}

func Test_gitPush_Success(t *testing.T) {
	// Open local repo and check out a test branch
	repo, err := openLocalRepo(localRepoDir)
	assert.NoError(t, err)
	worktree, err := repo.Worktree()
	assert.NoError(t, err)

	log.Info("git fetch")
	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{"refs/heads/integration-test:refs/heads/integration-test"},
		//Depth:    0,
		Progress: os.Stdout,
		//Tags:            0,
	})
	assert.NoError(t, err)

	assert.False(t, true)

	//log.Info("git pull")
	//err = worktree.Pull(&git.PullOptions{
	//	SingleBranch:      false,
	//	Depth:             0,
	//	RecurseSubmodules: 0,
	//	Progress:          os.Stdout,
	//})
	//assert.NoError(t, err)

	err = checkOutBranch(worktree, testBranch)
	assert.NoError(t, err)

	//log.Info("git pull")
	//err = worktree.Pull(&git.PullOptions{
	//	SingleBranch:      false,
	//	Depth:             0,
	//	RecurseSubmodules: 0,
	//	Progress:          os.Stdout,
	//})
	//assert.NoError(t, err)

	// Write a new file in the git repository
	err = ioutil.WriteFile(filepath.Join(localRepoDir, "foo-new-file.txt"), []byte("test123"), 0644)
	assert.NoError(t, err)

	// Git add and commit
	assert.NoError(t, err)
	assert.NoError(t, gitAddFiles([]string{"foo-new-file.txt"}, worktree))
	_, err = createCommit(worktree, "test commit", testUser, testEmail)
	assert.NoError(t, err)

	err = gitPush(repo, testUser, testGithubToken)
	assert.NoError(t, err)
}

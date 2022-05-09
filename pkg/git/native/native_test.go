package native

import (
	"github.com/champ-oss/file-sync/pkg/common"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const localRepoDir = "./"

var testBranch = "integration-test"
var testGithubToken = os.Getenv("GITHUB_TOKEN")
var testUser = "integration-test"
var testEmail = "integration-test@example.com"

func Test_cloneSourceRepo_Success(t *testing.T) {
	repoDir, err := cloneSourceRepo("https://github.com/git-fixtures/basic.git")
	defer common.RemoveDir(repoDir)
	assert.NoError(t, err)

	_, err = os.Stat(repoDir)
	assert.NoError(t, err)
}

func Test_cloneSourceRepo_Error(t *testing.T) {
	repoDir, err := cloneSourceRepo("https://localhost/not-a-repo")
	defer common.RemoveDir(repoDir)
	assert.Error(t, err)

	_, err = os.Stat(repoDir)
	assert.NoError(t, err)
}

func Test_openLocalRepo_Success(t *testing.T) {
	repoDir, err := cloneSourceRepo("https://github.com/git-fixtures/basic.git")
	defer common.RemoveDir(repoDir)

	repo, err := openLocalRepo(repoDir)
	assert.NoError(t, err)
	assert.NotNil(t, repo)
}

func Test_openLocalRepo_Error(t *testing.T) {
	repo, err := openLocalRepo(os.TempDir())
	assert.Error(t, err)
	assert.Nil(t, repo)
}

func Test_isWorktreeModified_Modified(t *testing.T) {
	// Clone and open an example git repository
	repoDir, err := cloneSourceRepo("https://github.com/git-fixtures/basic.git")
	defer common.RemoveDir(repoDir)
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
	defer common.RemoveDir(repoDir)
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
	defer common.RemoveDir(repoDir)
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
	defer common.RemoveDir(repoDir)
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
	defer common.RemoveDir(repoDir)
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
	defer common.RemoveDir(repoDir)
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
	defer common.RemoveDir(repoDir)
	repo, err := openLocalRepo(repoDir)
	assert.NoError(t, err)

	worktree, err := repo.Worktree()
	assert.NoError(t, err)
	assert.Error(t, gitAddFiles([]string{"foo-new-file.txt"}, worktree))
}

func Test_createCommit_Success(t *testing.T) {
	// Clone and open an example git repository
	repoDir, err := cloneSourceRepo("https://github.com/git-fixtures/basic.git")
	defer common.RemoveDir(repoDir)
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

//
//func Test_gitPush_Success(t *testing.T) {
//	// Open local repo and check out a test branch
//	repo, err := openLocalRepo(localRepoDir)
//	assert.NoError(t, err)
//	worktree, err := repo.Worktree()
//	assert.NoError(t, err)
//
//	log.Info("git fetch")
//	err = repo.Fetch(&git.FetchOptions{
//		//RemoteName: "origin",
//		RefSpecs: []config.RefSpec{"refs/heads/integration-test:refs/heads/integration-test"},
//		//Depth:    0,
//		Progress: os.Stdout,
//		//Tags:            0,
//	})
//	assert.NoError(t, err)
//
//	assert.False(t, true)
//
//	//log.Info("git pull")
//	//err = worktree.Pull(&git.PullOptions{
//	//	SingleBranch:      false,
//	//	Depth:             0,
//	//	RecurseSubmodules: 0,
//	//	Progress:          os.Stdout,
//	//})
//	//assert.NoError(t, err)
//
//	err = checkOutBranch(worktree, testBranch)
//	assert.NoError(t, err)
//
//	log.Info("git pull")
//	err = worktree.Pull(&git.PullOptions{
//		SingleBranch:      false,
//		Depth:             0,
//		RecurseSubmodules: 0,
//		Progress:          os.Stdout,
//	})
//	assert.NoError(t, err)
//
//	// Write a new file in the git repository
//	err = ioutil.WriteFile(filepath.Join(localRepoDir, "foo-new-file.txt"), []byte("test123"), 0644)
//	assert.NoError(t, err)
//
//	modified, err := isWorktreeModified(worktree)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if modified {
//		// Git add and commit
//		assert.NoError(t, err)
//		assert.NoError(t, gitAddFiles([]string{"foo-new-file.txt"}, worktree))
//		_, err = createCommit(worktree, "test commit", testUser, testEmail)
//		assert.NoError(t, err)
//
//		err = gitPush(repo, testUser, testGithubToken)
//		assert.NoError(t, err)
//	} else {
//		log.Info("no changes")
//	}
//
//}

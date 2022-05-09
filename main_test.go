package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

var testBranch = "integration-test"
var testGithubToken = os.Getenv("GITHUB_TOKEN")
var testUser = "integration-test"
var testEmail = "integration-test@example.com"

func init() {
	log.SetLevel(log.DebugLevel)
}

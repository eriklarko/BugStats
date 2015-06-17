package main

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"log"
	"asd"
	"encoding/json"
	"strings"
)

var session *sh.Session

func main() {
	session = sh.NewSession()
	session.SetDir("/home/erik/Code/NetClean/proactive")

	gitLog := getLog()
	for _, commit := range(gitLog) {
		if isBugFixCommit(&commit) {
			fmt.Printf("%s indicates a bugfix\n", commit.Message)
		}
	}
}

func getLog() []asd.HashAndMessage {
	rawOutput, err := session.Command("git", "log", "--pretty={\"hash\": \"%H\", \"message\": \"%f\"},").Output()
	if err != nil {
		log.Panicf("Unable to get the git log, %v\n", err)
	}

	padded := []byte("[" + string(rawOutput[:len(rawOutput)-2]) + "]")

	var gitLog []asd.HashAndMessage
	err = json.Unmarshal(padded, &gitLog)
	if err != nil {
		log.Panicf("Unable to parse the git log, %v\n", err)
	}
	return gitLog
}

func isBugFixCommit(commit *asd.HashAndMessage) bool {
	message := strings.ToLower(commit.Message)
	return strings.Contains(message, "fixes") || strings.Contains(message, "bugfix") || strings.Contains(message, "bug-fix")
}
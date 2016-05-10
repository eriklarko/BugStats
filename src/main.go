package main

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"log"
	"git"
	"strings"
	"methodnameextractor"
	"errors"
	"os"
)

func main() {
	git.MethodNameExtractors[".cs"] = methodnameextractor.GetMethodNamesFromLineCsharp
	git.MethodNameExtractors[".java"] = methodnameextractor.GetMethodNamesFromLineJava

	files, err := findFilesModified("/home/erik/Code/YetAnotherGTDApp")
	if err != nil {
		log.Panic(err)
	}

	outputFile, err := os.Create("./output")
	if err != nil {
		log.Panic(err)
	}

	fmt.Fprintf(outputFile, "file;count\n")
	for file, count := range files {
		fmt.Fprintf(outputFile, "%s;%d\n", file, count)
	}
	fmt.Printf("Wrote %d row(s) to %s\n", len(files), outputFile.Name())
}

func findFilesModified(pathToRepo string) (map[string]uint, error) {
	session := sh.NewSession()
	session.SetDir(pathToRepo)

	gitLog, err := git.GetCommitHistory(session)
	if err != nil {
		return nil, errors.New("Unable get the commit history. " + err.Error())
	}
	return findModifiedFiles(pathToRepo, gitLog), nil
}

func findFilesModifiedInBugfixCommits(pathToRepo string) (map[string]uint, error) {
	session := sh.NewSession()
	session.SetDir(pathToRepo)

	bugFixCommits, err := findBugFixCommits(pathToRepo)
	if err != nil {
		return nil, err
	}
	return findModifiedFiles(pathToRepo, bugFixCommits), nil
}

func findModifiedFiles(pathToRepo string, commits []*git.HashAndMessage) map[string]uint {
	session := sh.NewSession()
	session.SetDir(pathToRepo)

	fileCounts := make(map[string]uint)
	for _, commit := range(commits) {
		fileChanges, err := git.GetModifiedFiles(session, commit.Hash)
		if err != nil {
			log.Printf("Unable to parse %+v, %+v\n", commit, err)
		}
		for _, fileChange := range fileChanges {
			fileName := fileChange.FileName
			if _, found := fileCounts[fileName]; found {
				fileCounts[fileName]++
			} else {
				fileCounts[fileName] = 1
			}
		}
	}
	return fileCounts
}

func findBugFixCommits(pathToRepo string) ([]*git.HashAndMessage, error) {
	session := sh.NewSession()
	session.SetDir(pathToRepo)

	gitLog, err := git.GetCommitHistory(session)
	if err != nil {
		return nil, errors.New("Unable get the commit history. " + err.Error())
	}

	bugFixCommits := make([]*git.HashAndMessage, 0)
	for _, commit := range(gitLog) {
		if isBugFixCommit(commit) {
			fmt.Printf("%s indicates a bugfix (%s)\n", commit.Message, commit.Hash)
			bugFixCommits = append(bugFixCommits, commit)
		}
	}
	return bugFixCommits, nil
}

func findModifiedMethods(pathToRepo string, commits []*git.HashAndMessage) map[string]uint {
	session := sh.NewSession()
	session.SetDir(pathToRepo)

	methodCounts := make(map[string]uint)
	for _, commit := range(commits) {
		methods := git.GetModifiedMethods(session, commit)

		for _, method := range methods {
			if _, found := methodCounts[method]; found {
				methodCounts[method]++
			} else {
				methodCounts[method] = 1
			}
		}
	}
	return methodCounts
}

func findMethodsModifiedInCommit(pathToRepo string, commit *git.HashAndMessage) []string {
	session := sh.NewSession()
	session.SetDir(pathToRepo)

	return git.GetModifiedMethods(session, commit);
}

func isBugFixCommit(commit *git.HashAndMessage) bool {
	message := strings.ToLower(commit.Message)
	return strings.Contains(message, "fixes") || strings.Contains(message, "bugfix") || strings.Contains(message, "bug-fix")
}

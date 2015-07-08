package main

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"log"
	"git"
	"strings"
	"methodnameextractor"
	"errors"
	"path/filepath"
)

var session *sh.Session
var methodNameExtractors map[string]func(string, uint) string = make(map[string]func(string, uint) string)

func main() {
	methodNameExtractors[".cs"] = methodnameextractor.GetMethodNameFromLineCsharp
	methodNameExtractors[".java"] = methodnameextractor.GetMethodNameFromLineJava
	
	session = sh.NewSession()
	session.SetDir("/home/erik/Code/NetClean/proactive")


	analyze(&git.HashAndMessage{Hash: "f2da0220bbaf29afb769df64e230dc4c4828d2bf"})
	log.Fatalln("BYE!") // TODO: Remove

	gitLog := git.GetLog(session)
	for _, commit := range(gitLog) {
		if isBugFixCommit(&commit) {
			fmt.Printf("%s indicates a bugfix (%s)\n", commit.Message, commit.Hash)
			analyze(&commit)
		}
	}
}

func isBugFixCommit(commit *git.HashAndMessage) bool {
	message := strings.ToLower(commit.Message)
	return strings.Contains(message, "fixes") || strings.Contains(message, "bugfix") || strings.Contains(message, "bug-fix")
}

func analyze(commit *git.HashAndMessage) {
	modifiedFiles := git.GetModifiedFiles(session, commit.Hash)
	fmt.Printf("  Modified files: %v\n", len(modifiedFiles))
	for _, modifiedFile := range(modifiedFiles) {
		if (modifiedFile.FileChange != git.MODIFIED) {
			continue
		}

		modifiedLines := git.GetLinesModifiedInFile(session, commit.Hash, modifiedFile.FileName)
		fmt.Printf("    Lines in %s, %v\n", modifiedFile.FileName, modifiedLines)

		contents, err := git.GetFileContents(session, commit.Hash, modifiedFile.FileName)
		if err != nil {
			log.Panicln(err.Error())
		}

		for _, modifiedLine := range modifiedLines {
			methodName, err := getMethodNameFromLine(modifiedFile.FileName, contents, modifiedLine)
			if err == nil {
				log.Printf("The method name on %s:%d is %s\n", modifiedFile.FileName, modifiedLine, methodName)
				// TODO: Use this method name :) Print it to CSV or something.
			} else {
				log.Printf("Unable to get method name from %s:%d - %v\n", modifiedFile.FileName, modifiedLine, err)
			}
		}
	}
}

func getMethodNameFromLine(fileName string, fileContents string, lineNumber uint) (string, error) {
	//fmt.Printf("Getting the method name on line %d in %s\n",lineNumber, fileName)

	fileEnding := filepath.Ext(fileName)
	extractor := methodNameExtractors[fileEnding]
	if extractor == nil {
		return "", errors.New("No method name extractor found for " + fileEnding)
	}

	return extractor(fileContents, lineNumber), nil
}
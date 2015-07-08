package main

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"log"
	"git"
	"strings"
	"methodnameextractor"
	"path/filepath"
)

var session *sh.Session
var methodNameExtractors map[string]func(string, uint) (string, error) = make(map[string]func(string, uint) (string, error))

func main() {
	methodNameExtractors[".cs"] = methodnameextractor.GetMethodNameFromLineCsharp
	methodNameExtractors[".java"] = methodnameextractor.GetMethodNameFromLineJava

	session = sh.NewSession()
	session.SetDir("/home/erik/Code/NetClean/proactive")


	analyzeBuggyMethods(&git.HashAndMessage{Hash: "f2da0220bbaf29afb769df64e230dc4c4828d2bf"})
	log.Fatalln("BYE!") // TODO: Remove

	gitLog, err := git.GetLog(session)
	if err != nil {
		log.Fatalln(err)
	}
	for _, commit := range(gitLog) {
		if isBugFixCommit(&commit) {
			fmt.Printf("%s indicates a bugfix (%s)\n", commit.Message, commit.Hash)
			analyzeBuggyMethods(&commit)
		}
	}
}

func isBugFixCommit(commit *git.HashAndMessage) bool {
	message := strings.ToLower(commit.Message)
	return strings.Contains(message, "fixes") || strings.Contains(message, "bugfix") || strings.Contains(message, "bug-fix")
}

func analyzeBuggyMethods(commit *git.HashAndMessage) {
	modifiedFiles, err := git.GetModifiedFiles(session, commit.Hash)
	if err != nil {
		log.Panicln(err)
	}

	for _, modifiedFile := range(modifiedFiles) {
		if (modifiedFile.FileChange != git.MODIFIED) {
			continue
		}

		modifiedLines, err := git.GetLinesModifiedInFile(session, commit.Hash, modifiedFile.FileName)
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Printf("  Lines in %s, %v\n", modifiedFile.FileName, modifiedLines)

		contents, err := git.GetFileContents(session, commit.Hash, modifiedFile.FileName)
		if err != nil {
			log.Println(err)
			continue
		}


		fileEnding := filepath.Ext(modifiedFile.FileName)
		extractor := methodNameExtractors[fileEnding]
		if extractor == nil {
			log.Println("No method name extractor found for " + fileEnding)
			continue
		}

		for _, modifiedLine := range modifiedLines {
			// TODO: Invoke the secret sauce java program.
			//fmt.Printf("Getting the method name on line %d in %s\n",lineNumber, fileName)
			methodName, err := extractor(contents, modifiedLine)

			if err == nil {
				log.Printf("The method name on %s:%d is %s\n", modifiedFile.FileName, modifiedLine, methodName)
				// TODO: Use this method name :) Print it to CSV or something.
			} else {
				log.Printf("Unable to get method name from %s:%d - %v\n", modifiedFile.FileName, modifiedLine, err)
			}
		}
	}
}
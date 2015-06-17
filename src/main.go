package main

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"log"
	"git"
	"encoding/json"
	"strings"
	"regexp"
	"strconv"
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

	gitLog := getLog()[6:] // TODO: Remove skip here
	for _, commit := range(gitLog) {
		if isBugFixCommit(&commit) {
			fmt.Printf("%s indicates a bugfix\n", commit.Message)

			modifiedFiles := getModifiedFiles(commit.Hash)
			for _, modifiedFile := range(modifiedFiles) {
				modifiedLines := getLinesModifiedInFile(commit.Hash, modifiedFile)

				for _, modifiedLine := range modifiedLines {
					contents := getFileContents(commit.Hash, modifiedFile)
					methodName, err := getMethodNameFromLine(modifiedFile, contents, modifiedLine)
					if err == nil {
						log.Printf("The method name on %s:%d is %s\n", modifiedFile, modifiedLine, methodName)
						// TODO: Use this method name :) Print it to CSV or something.
					} else {
						log.Printf("Unable to get method name from %s:%d - %v\n", modifiedFile, modifiedLine, err)
					}
				}

				log.Fatalln("BYE!") // TODO: Remove
			}
		}
	}
}

func getLog() []git.HashAndMessage {
	rawOutput, err := session.Command("git", "log", "--pretty={\"hash\": \"%H\", \"message\": \"%f\"},").Output()
	if err != nil {
		log.Panicf("Unable to get the git log, %v\n", err)
	}

	padded := []byte("[" + string(rawOutput[:len(rawOutput)-2]) + "]")

	var gitLog []git.HashAndMessage
	err = json.Unmarshal(padded, &gitLog)
	if err != nil {
		log.Panicf("Unable to parse the git log, %v\n", err)
	}
	return gitLog
}

func isBugFixCommit(commit *git.HashAndMessage) bool {
	message := strings.ToLower(commit.Message)
	return strings.Contains(message, "fixes") || strings.Contains(message, "bugfix") || strings.Contains(message, "bug-fix")
}

// Returns a list of file names of all files modified in the commit and it's first parent
// A merged feature branch's first parent is develop (or whichever branch it was merged into)
func getModifiedFiles(commitHash string) []string {
	// TODO: How does the git diff --name-only output look for moved files?
	// TODO: Does this work for the first commit?

	rawOutput, err := session.Command("git", "diff", "--name-only", commitHash, commitHash + "^").Output()
	if err != nil {
		log.Panicf("Unable to get list of files modified in %s. %v\n", commitHash, err)
	}
	return strings.Split(string(rawOutput), "\n")
}

// TODO: Not tested enough
func getLinesModifiedInFile(commitHash string, file string) []uint {
	// git diff commitHash commitHash^ -- "file"
	cmd := session.Command("bash", "-c", "git diff -U0 " + commitHash + " " + commitHash + "^ -- \"" + file + "\"")
	cmd.ShowCMD = true
	rawOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Panicf("Unable to get which lines was modified in %s.\n%s\n", file, rawOutput)
	}

	unifiedDiffAffectedLinesRegExp := regexp.MustCompile("(@@.*@@)")
	affectedLinesData := unifiedDiffAffectedLinesRegExp.FindAllString(string(rawOutput), -1)

	totalAffectedRows := make(map[uint]struct{})
	unifiedDiffLinesChangedInOldFileRegExp := regexp.MustCompile("-(\\d+)(,(\\d+))?")
	for _, affectedLinesDatum := range (affectedLinesData) {
		raw := unifiedDiffLinesChangedInOldFileRegExp.FindAllStringSubmatch(affectedLinesDatum, -1)
		if len(raw) != 1 {
			log.Panicf("Something went wrong parsing %s, got wrong number of outer groups (%v)\n", affectedLinesDatum, raw)
		}

		var rawRow string
		var rawNumberOfRows string
		if len(raw[0]) == 4 {
			rawRow = raw[0][1]
			rawNumberOfRows = raw[0][3]
			if len(rawNumberOfRows) == 0 {
				rawNumberOfRows = "1"
			}
		} else {
			log.Panicf("Something went wrong parsing %s, got wrong number of inner groups (%v)\n", affectedLinesDatum, raw)
		}

		row, err := strconv.Atoi(rawRow)
		if err != nil {
			log.Panicf("Something went wrong parsing %s, the row is not a number (%v)\n", affectedLinesDatum, rawRow)
		}
		numRows, err := strconv.Atoi(rawNumberOfRows)
		if err != nil {
			log.Panicf("Something went wrong parsing %s, the number of rows is not a number (%v)\n", affectedLinesDatum, rawNumberOfRows)
		}

		affectedRows := expandRowNumberAndNumberOfAffectedRows(uint(row), uint(numRows))
		addAll(affectedRows, totalAffectedRows)
	}

	return keys(totalAffectedRows)
}

func expandRowNumberAndNumberOfAffectedRows(row uint, numberOfRows uint) []uint {
	rows := make([]uint, numberOfRows)
	for i := uint(0); i < numberOfRows; i++ {
		rows[i] = row + i
	}
	return rows
}

func addAll(toAdd []uint, target map[uint]struct{}) {
	var a struct{}
	for _, k := range toAdd {
		target[k] = a
	}
}

func keys (theMap map[uint]struct{}) []uint {
	keys := make([]uint, 0, len(theMap))
	for k := range theMap {
		keys = append(keys, k)
	}
	return keys
}

func getFileContents(commitHash string, file string) string {
	// git show hash:"file"
	rawOutput, err := session.Command("git", "show", commitHash+":"+file).CombinedOutput()
	if err != nil {
		log.Panicf("Unable to get which lines was modified in %s.\n%s\n", file, rawOutput)
	}
	return string(rawOutput)
}

func getMethodNameFromLine(fileName string, fileContents string, lineNumber uint) (string, error) {
	fmt.Printf("Getting the method name on line %d in %s\n",lineNumber, fileName)

	fileEnding := filepath.Ext(fileName)
	extractor := methodNameExtractors[fileEnding]
	if extractor == nil {
		return "", errors.New("No method name extractor found for " + fileEnding)
	}

	return extractor(fileContents, lineNumber), nil
}
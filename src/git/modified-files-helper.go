package git

import (
	"log"
	"errors"
	"strings"
	"github.com/codeskyblue/go-sh"
)

// Returns a list of file names of all files modified in the commit and it's first parent
// A merged feature branch's first parent is develop (or whichever branch it was merged into)
func GetModifiedFiles(session *sh.Session, commitHash string) []*FileChange {
	// TODO: How does the git diff --name-only output look for moved files?
	// TODO: Does this work for the first commit?

	rawOutput, err := session.Command("git", "diff", "--name-status", commitHash, commitHash + "^").Output()
	if err != nil {
		log.Panicf("Unable to get list of files modified in %s. %v\n", commitHash, err)
	}

	rawRows := strings.Split(string(rawOutput), "\n")
	fileChanges := make([]*FileChange, 0)
	for _, rawRow := range(rawRows) {
		if (len(rawRow) == 0) {
			continue
		}

		change, err := getChangeFromChar(rawRow[0])
		if err != nil {
			log.Panicf("Unable to change to file %s. %v\n", rawRow, err)
		}
		lineParts := strings.Split(rawRow, "\t")
		name := lineParts[1]
		fileChanges = append(fileChanges, &FileChange{FileChange: change, FileName: name})
	}
	return fileChanges
}

func getChangeFromChar(c uint8) (Change, error) {
	switch string(c) {
	case "A":
		return CREATED, nil
	case "M":
		return MODIFIED, nil
	case "D":
		return REMOVED, nil
	default:
		return -1, errors.New("Unknown change: " + string(c))
	}
}
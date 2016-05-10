package git

import (
	"strconv"
	"regexp"
	"github.com/codeskyblue/go-sh"
	"errors"
	"fmt"
)

// TODO: Not tested enough
func GetLinesModifiedInFile(session *sh.Session, commitHash string, file string) ([]uint, error) {
	// git diff commitHash commitHash^ -- "file"
	cmd := session.Command("bash", "-c", "git diff -U0 " + commitHash + " " + commitHash + "^ -- \"" + file + "\"")
	rawOutput, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to get which lines was modified in %s.\n%s\n", file, rawOutput))
	}

	unifiedDiffAffectedLinesRegExp := regexp.MustCompile("(@@.*@@)")
	affectedLinesData := unifiedDiffAffectedLinesRegExp.FindAllString(string(rawOutput), -1)

	totalAffectedRows := make(map[uint]struct{})
	for _, hunkHead := range (affectedLinesData) {
		affectedRows, err := getRowsFromHunkHead(hunkHead)
		if err != nil {
			return nil, err
		}
		addAll(affectedRows, totalAffectedRows)
	}

	return keys(totalAffectedRows), nil
}

func getRowsFromHunkHead(hunkHead string) ([]uint, error) {
	unifiedDiffLinesChangedInOldFileRegExp := regexp.MustCompile("-(\\d+)(,(\\d+))? ")
	raw := unifiedDiffLinesChangedInOldFileRegExp.FindAllStringSubmatch(hunkHead, -1)
	if len(raw) != 1 {
		return nil, errors.New(fmt.Sprintf("Something went wrong parsing %s, got wrong number of outer groups (%v)\n", hunkHead, raw))
	}

	var rawRow string
	var rawNumberOfRows string
	if len(raw[0]) == 4 {
		rawRow = raw[0][1]
		rawNumberOfRows = raw[0][3]
		if len(rawNumberOfRows) == 0 {
			rawNumberOfRows = "0"
		}
	} else {
		return nil, errors.New(fmt.Sprintf("Something went wrong parsing %s, got wrong number of inner groups (%v)\n", hunkHead, raw))
	}

	row, err := strconv.Atoi(rawRow)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Something went wrong parsing %s, the row is not a number (%v)\n", hunkHead, rawRow))
	}
	numRows, err := strconv.Atoi(rawNumberOfRows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Something went wrong parsing %s, the number of rows is not a number (%v)\n", hunkHead, rawNumberOfRows))
	}

	//fmt.Printf("    %s -> %v, %v\n", hunkHead, row, numRows)
	return expandRowNumberAndNumberOfAffectedRows(uint(row), uint(numRows)), nil
}

func expandRowNumberAndNumberOfAffectedRows(row uint, numberOfRows uint) []uint {
	rows := make([]uint, numberOfRows + 1)
	for i := uint(0); i < uint(len(rows)); i++ {
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

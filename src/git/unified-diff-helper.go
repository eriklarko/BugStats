package git

import (
	"log"
	"strconv"
	"regexp"
	"github.com/codeskyblue/go-sh"
)

// TODO: Not tested enough
func GetLinesModifiedInFile(session *sh.Session, commitHash string, file string) []uint {
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
	for _, hunkHead := range (affectedLinesData) {
		affectedRows := getRowsFromHunkHead(hunkHead)
		addAll(affectedRows, totalAffectedRows)
	}

	return keys(totalAffectedRows)
}

func getRowsFromHunkHead(hunkHead string) []uint {
	unifiedDiffLinesChangedInOldFileRegExp := regexp.MustCompile("-(\\d+)(,(\\d+))?")
	raw := unifiedDiffLinesChangedInOldFileRegExp.FindAllStringSubmatch(hunkHead, -1)
	if len(raw) != 1 {
		log.Panicf("Something went wrong parsing %s, got wrong number of outer groups (%v)\n", hunkHead, raw)
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
		log.Panicf("Something went wrong parsing %s, got wrong number of inner groups (%v)\n", hunkHead, raw)
	}

	row, err := strconv.Atoi(rawRow)
	if err != nil {
		log.Panicf("Something went wrong parsing %s, the row is not a number (%v)\n", hunkHead, rawRow)
	}
	numRows, err := strconv.Atoi(rawNumberOfRows)
	if err != nil {
		log.Panicf("Something went wrong parsing %s, the number of rows is not a number (%v)\n", hunkHead, rawNumberOfRows)
	}

	//fmt.Printf("    %s -> %v, %v\n", hunkHead, row, numRows)
	return expandRowNumberAndNumberOfAffectedRows(uint(row), uint(numRows))
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

package git

import (
	"log"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codeskyblue/go-sh"
)

func GetLog(session *sh.Session) []HashAndMessage {
	rawOutput, err := session.Command("git", "log", "--pretty={\"hash\": \"%H\", \"message\": \"%f\"},").Output()
	if err != nil {
		log.Panicf("Unable to get the git log, %v\n", err)
	}

	padded := []byte("[" + string(rawOutput[:len(rawOutput)-2]) + "]")

	var gitLog []HashAndMessage
	err = json.Unmarshal(padded, &gitLog)
	if err != nil {
		log.Panicf("Unable to parse the git log, %v\n", err)
	}
	return gitLog
}


func GetFileContents(session *sh.Session, commitHash string, file string) (string, error) {
	// git show hash:"file"
	rawOutput, err := session.Command("git", "show", commitHash+":"+file+"").CombinedOutput()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to get which lines was modified in %s.\n%s\n", file, rawOutput))
	}
	return string(rawOutput), nil
}
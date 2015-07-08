package git

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codeskyblue/go-sh"
)

func GetLog(session *sh.Session) ([]HashAndMessage, error) {
	rawOutput, err := session.Command("git", "log", "--pretty={\"hash\": \"%H\", \"message\": \"%f\"},").Output()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to get the git log, %v", err))
	}

	padded := []byte("[" + string(rawOutput[:len(rawOutput)-2]) + "]")

	var gitLog []HashAndMessage
	err = json.Unmarshal(padded, &gitLog)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to parse the git log, %v", err))
	}
	return gitLog, nil
}


func GetFileContents(session *sh.Session, commitHash string, file string) (string, error) {
	// git show hash:"file"
	rawOutput, err := session.Command("git", "show", commitHash+":"+file+"").CombinedOutput()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to get which lines was modified in %s.\n%s\n", file, rawOutput))
	}
	return string(rawOutput), nil
}
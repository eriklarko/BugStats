package main

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"log"
	"asd"
	"encoding/json"
)

var session *sh.Session

func main() {
	// git log --pretty="{\"hash\": \"%H\", \"message\": \"%B\"}"
	session = sh.NewSession()
	session.SetDir("/home/erik/Code/Privat/BrewersLittleHelper")

	gitLog := getLog()
	fmt.Printf("%v", gitLog)
}

func getLog() []asd.HashAndMessage {
	rawOuput, err := session.Command("git", "log", "--pretty={\"hash\": \"%H\", \"message\": \"%f\"},").Output()
	if err != nil {
		log.Panicf("Unable to get the git log, %v\n", err)
	}

	padded := []byte("[" + string(rawOuput[:len(rawOuput)-2]) + "]")

	var gitLog []asd.HashAndMessage
	err = json.Unmarshal(padded, &gitLog)
	if err != nil {
		log.Panicf("Unable to parse the git log, %v\n", err)
	}
	return gitLog
}
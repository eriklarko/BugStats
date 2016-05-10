package git

import (
	"log"
	"path/filepath"
	"github.com/codeskyblue/go-sh"
	"errors"
	"io/ioutil"
	"set"
)

var MethodNameExtractors map[string]func(string, []uint) ([]string, error) = make(map[string]func(string, []uint) ([]string, error))


func GetModifiedMethods(session *sh.Session, commit *HashAndMessage) []string {
	modifiedFiles, err := GetModifiedFiles(session, commit.Hash)
	if err != nil {
		log.Panicln(err)
	}

	allModifiedMethods := set.NewSet()
	for _, modifiedFile := range(modifiedFiles) {
		if (modifiedFile.FileChange != MODIFIED) {
			continue
		}

		modifiedMethods, err := findModifiedMethodsInFile(session, commit, modifiedFile);
		if err == nil {
			allModifiedMethods.AddAll(modifiedMethods)
		} else {
			log.Printf("Failed getting modified methods from file %s, %v\n", modifiedFile.FileName, err)
		}
	}

	return allModifiedMethods.AsSlice()
}

func findModifiedMethodsInFile(session *sh.Session, commit *HashAndMessage, modifiedFile *FileChange) ([]string, error) {
	modifiedLines, err := GetLinesModifiedInFile(session, commit.Hash, modifiedFile.FileName)
	if err != nil {
		return nil, err
	}

	fileEnding := filepath.Ext(modifiedFile.FileName)
	extractor := MethodNameExtractors[fileEnding]
	if extractor == nil {
		return nil, errors.New("No method name extractor found for " + fileEnding)
	}

	contents, err := GetFileContents(session, commit.Hash, modifiedFile.FileName)
	if err != nil {
		return nil, errors.New("Unable to get the contents of " + modifiedFile.FileName + ". " + err.Error())
	}

	err = ioutil.WriteFile("/tmp/apa", []byte(contents), 0644)
	if err != nil {
		return nil, errors.New("Unable to write contents to /tmp/apa, " + err.Error())
	}

	modifiedMethods, err := extractor("/tmp/apa", modifiedLines)
	if err != nil {
		return nil, err
	}

	for i, modifiedMethod := range modifiedMethods {
		modifiedMethods[i] = modifiedFile.FileName + ":" + modifiedMethod
	}
	return modifiedMethods, nil
}


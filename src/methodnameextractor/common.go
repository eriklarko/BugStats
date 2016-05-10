package methodnameextractor

import (
	"strconv"
	"encoding/json"
	"errors"
)

func linesToString(lines []uint) string {
	if len(lines) == 0 {
		return ""
	}

	str := ""
	for _, line := range lines {
		str += strconv.Itoa(int(line)) + " "
	}
	return str[:len(str) - 1]
}
/**
 * Expects the byte slice from the output of the magic sauce program that
 * determines which if any method is on a line. The output should be on the
 * form
    [
	{"line":53,"hasMethod":true,"methodName":"handleIncomingFile","parameterTypes":["Request","Response"]},
 	{"line":90,"hasMethod":false,"methodName":"","parameterTypes":[]},
    ]

 */
func getMethodNamesFromRawOutput(raw []byte) ([]string, error) {
	var parsed []methodNameExtractorOutput
	err := json.Unmarshal(raw, &parsed)
	if err != nil {
		return nil, errors.New("Unable to parse the output as json, " + err.Error());
	}

	names := make([]string, 0)
	for _, parsedMethod := range parsed {
		if parsedMethod.HasMethod {
			names = append(names, parsedMethod.MethodName)
		}
	}
	return names, nil;
}


type methodNameExtractorOutput struct {
	Line int `json:"line"`
	HasMethod bool `json:"hasMethod"`
	MethodName string `json:"methodName"`
	ParameterTypes []string `json:"parameterTypes"`
}

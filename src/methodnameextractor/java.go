package methodnameextractor

import (
  "os/exec"
  "fmt"
  "errors"
)

func GetMethodNamesFromLineJava(file string, lines []uint) ([]string, error) {
  // java -jar target/methodnameextractor-1.0-SNAPSHOT-jar-with-dependencies.jar file line

  jarFilePath := "method-name-extractor/target/methodnameextractor-1.0-SNAPSHOT-jar-with-dependencies.jar"

  cmd := exec.Command("java", "-jar", jarFilePath, file, linesToString(lines))
  rawOutput, err := cmd.CombinedOutput()
  if err != nil {
    fmt.Println(string(rawOutput))
    return nil, errors.New("Failed invoking magic sauce java program, " + err.Error())
  }

  return getMethodNamesFromRawOutput(rawOutput)
}


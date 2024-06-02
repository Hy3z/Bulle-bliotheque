package util

import (
	"bytes"
	"os"
)

const (
	CypherScriptDirectory = "./script"
)

func ReadCypherScript(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	script := []byte{}
	for _, char := range data {
		if !bytes.Equal([]byte{char}, []byte("\n")) {
			script = append(script, char)
		} else {
			script = append(script, ' ')
		}
	}

	return string(script), nil
}

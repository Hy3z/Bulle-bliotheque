package util

import (
	"bytes"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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

func RecordsContains(records []*neo4j.Record, index int, value any) bool {
	for _, record := range records {
		v := record.Values[index]
		if v == value {
			return true
		}
		continue
	}
	return false
}

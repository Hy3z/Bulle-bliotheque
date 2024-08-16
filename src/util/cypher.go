package util

import (
	"bb/database"
	"bytes"
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"os"
)

const (
	//Dossier où sont situés les requêtes cypher
	CypherScriptDirectory = "./script"
)

// ReadCypherScript lit un fichier et en renvoit son contenu sous forme de string
func ReadCypherScript(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	var script []byte
	for _, char := range data {
		if !bytes.Equal([]byte{char}, []byte("\n")) {
			script = append(script, char)
		} else {
			script = append(script, ' ')
		}
	}

	return string(script), nil
}

// RecordsContains renvoit true si au moins une des réponse contient la valeur voulue à un certain index
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

// ExecuteCypherScript éxecute un script cypher contre la base de données et renvoie le résultat
func ExecuteCypherScript(filepath string, parameters map[string]any) (*neo4j.EagerResult, error) {
	query, err := ReadCypherScript(filepath)
	if err != nil {
		return nil, err
	}
	res, err := database.Query(context.Background(), query, parameters)
	return res, err
}

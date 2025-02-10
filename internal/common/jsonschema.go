package common

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

func LoadJSONSchema(filePath string) (*gojsonschema.Schema, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open schema file: %w", err)
	}
	defer file.Close()

	var schemaData interface{}
	if err := json.NewDecoder(file).Decode(&schemaData); err != nil {
		return nil, fmt.Errorf("failed to decode JSON schema: %w", err)
	}

	loader := gojsonschema.NewGoLoader(schemaData)
	return gojsonschema.NewSchema(loader)
}

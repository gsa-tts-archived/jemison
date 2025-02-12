package common

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/xeipuuv/gojsonschema"
)

// Embed the schemas directory.
//
//go:embed schemas/*.json
var schemasFS embed.FS

// LoadEmbeddedSchema loads a JSON schema from the embedded filesystem.
func LoadEmbeddedSchema(schemaName string) (*gojsonschema.Schema, error) {
	// Read the file contents from the embedded filesystem
	data, err := fs.ReadFile(schemasFS, "schemas/"+schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded schema %s: %w", schemaName, err)
	}

	// Create a loader from the schema content
	loader := gojsonschema.NewStringLoader(string(data))

	// Parse the schema and wrap any errors returned
	schema, err := gojsonschema.NewSchema(loader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema %s: %w", schemaName, err)
	}

	return schema, nil
}

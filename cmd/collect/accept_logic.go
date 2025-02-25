package main

import (
	"errors"
	"fmt"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/zap"
)

var (
	entreeSchema *gojsonschema.Schema
	fetchSchema  *gojsonschema.Schema
	initialized  = false
)

// InitializeSchemas initializes all necessary schemas (called on startup).
func InitializeSchemas() error {
	if initialized {
		return nil
	}

	var err error
	// Load schemas through embedded filesystem
	entreeSchema, err = common.LoadEmbeddedSchema("entree_schema.json")
	if err != nil {
		return fmt.Errorf("failed to load entree schema: %w", err)
	}

	fetchSchema, err = common.LoadEmbeddedSchema("fetch_schema.json")
	if err != nil {
		return fmt.Errorf("failed to load fetch schema: %w", err)
	}

	// Mark schemas as initialized
	initialized = true

	return nil
}

// ValidateJSON validates a JSON object against a schema.
func ValidateJSON(schema *gojsonschema.Schema, rawJSON string) error {
	zap.L().Info("Validating JSON", zap.String("json", rawJSON))
	documentLoader := gojsonschema.NewStringLoader(rawJSON)

	result, err := schema.Validate(documentLoader)
	if err != nil {
		return fmt.Errorf("schema validation error: %w", err)
	}

	if !result.Valid() {
		for _, desc := range result.Errors() {
			zap.L().Error("JSON validation error", zap.String("field", desc.Field()), zap.String("description", desc.Description())) //nolint:lll
		}

		return errors.New("JSON validation failed")
	}

	return nil
}

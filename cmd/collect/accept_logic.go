package main

import (
	"encoding/json"
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
	entreeSchema, err = common.LoadJSONSchema("cmd/collect/schemas/entree_schema.json")
	if err != nil {
		return fmt.Errorf("failed to load entree schema: %w", err)
	}

	fetchSchema, err = common.LoadJSONSchema("cmd/collect/schemas/fetch_schema.json")
	if err != nil {
		return fmt.Errorf("failed to load fetch schema: %w", err)
	}

	initialized = true
	return nil
}

// ValidateJSON validates a JSON object against a schema.
func ValidateJSON(schema *gojsonschema.Schema, rawJSON string) error {
	documentLoader := gojsonschema.NewStringLoader(rawJSON)
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return fmt.Errorf("schema validation error: %w", err)
	}

	if !result.Valid() {
		for _, desc := range result.Errors() {
			zap.L().Error("JSON validation error", zap.String("field", desc.Field()), zap.String("description", desc.Description()))
		}
		return errors.New("JSON validation failed")
	}
	return nil
}

func HandleBusinessLogic(args common.CollectArgs, jsonString string) error {
	zap.L().Info("Handling business logic", zap.String("json", jsonString))

	// Parse JSON into a map
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &jsonData); err != nil {
		zap.L().Error("JSON unmarshaling failed", zap.Error(err))
		return err
	}

	// Extract and validate `source` field from JSON
	source, ok := jsonData["source"].(string)
	if !ok || source == "" {
		return errors.New("missing or invalid `source` field in JSON")
	}

	// Log validation success
	zap.L().Info("JSON successfully validated and processed",
		zap.String("source", source))

	// Business logic based on `source`
	zap.L().Info("Processing source", zap.String("source", source))
	// Add source-specific processing logic here...

	return nil
}

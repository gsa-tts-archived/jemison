package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/riverqueue/river"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/zap"
)

//nolint:revive
func (w *CollectWorker) Work(ctx context.Context, job *river.Job[common.CollectArgs]) error {
	zap.L().Info("Collect Worker processing job", zap.String("data-json", job.Args.JSON))

	// Ensure schemas are initialized
	if err := ensureSchemasInitialized(); err != nil {
		return err
	}

	// Deserialize JSON data
	rawData, err := deserializeJSON(job.Args.JSON)
	if err != nil {
		return err
	}

	// Check for `source` field and select schema
	schema, err := selectSchema(rawData)
	if err != nil {
		return err
	}

	// Validate JSON and handle business logic
	if err := validateAndHandleBusinessLogic(schema, job.Args); err != nil {
		return err
	}

	return nil
}

func ensureSchemasInitialized() error {
	if !initialized {
		err := errors.New("schemas not initialized")
		zap.L().Fatal("schemas not initialized", zap.Error(err))

		return err
	}

	return nil
}

func deserializeJSON(jsonString string) (map[string]interface{}, error) {
	var rawData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &rawData); err != nil {
		// Log the error
		zap.L().Error("failed to unmarshal JSON", zap.Error(err))
		// Wrap the error with additional context and return
		return nil, fmt.Errorf("deserializeJSON: failed to unmarshal input JSON: %w", err)
	}

	// Pull in fullCrawl and hallpass for debugging or processing
	fullCrawl, _ := rawData["fullCrawl"].(bool)
	hallPass, _ := rawData["hallpass"].(bool)

	zap.L().Debug("deserialized JSON attributes",
		zap.Bool("fullCrawl", fullCrawl),
		zap.Bool("hallPass", hallPass),
		zap.Any("rawData", rawData))

	return rawData, nil
}

func selectSchema(rawData map[string]interface{}) (*gojsonschema.Schema, error) {
	// Extract the "data" object
	data, ok := rawData["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("missing or invalid `data` object")
	}

	// Extract the "source" field from the "data" object
	source, ok := data["source"].(string)
	if !ok || source == "" {
		return nil, errors.New("missing or invalid `source` field")
	}

	// Determine the schema based on the "source" value
	switch source {
	case "entree":
		return entreeSchema, nil
	case "fetch":
		return fetchSchema, nil
	default:
		return nil, errors.New("unsupported `source` value")
	}
}

func validateAndHandleBusinessLogic(schema *gojsonschema.Schema, args common.CollectArgs) error {
	if err := ValidateJSON(schema, args.JSON); err != nil {
		zap.L().Error("JSON validation failed", zap.Error(err))

		return err
	}

	// Log validation success
	zap.L().Info("JSON validation passed",
		zap.String("json", args.JSON))

	return nil
}

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

	// Deserialize JSON string into a map
	jsonData, err := deserializeJSON(job.Args.JSON)
	if err != nil {
		return err
	}

	// Check for `source` field and select schema
	jsonSchema, err := selectSchema(jsonData)
	if err != nil {
		return err
	}

	// Validate JSON
	if err := validateJSONData(jsonSchema, job.Args.JSON); err != nil {
		return err
	}

	// Send the validated data to S3
	if err := sendToS3(job.Args.JSON); err != nil {
		zap.L().Error("Failed to send data to S3", zap.Error(err))

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
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &jsonData); err != nil {
		zap.L().Error("failed to unmarshal JSON", zap.Error(err))

		return nil, fmt.Errorf("deserializeJSON: failed to unmarshal input JSON: %w", err)
	}

	// Pull in IsFull and hallpass for debugging or processing
	isFull, _ := jsonData["IsFull"].(bool)
	hallPass, _ := jsonData["hallpass"].(bool)

	zap.L().Debug("deserialized JSON attributes",
		zap.Bool("isFull", isFull),
		zap.Bool("hallPass", hallPass),
		zap.Any("jsonData", jsonData))

	return jsonData, nil
}

func selectSchema(jsonData map[string]interface{}) (*gojsonschema.Schema, error) {
	// Extract the "data" object
	data, ok := jsonData["data"].(map[string]interface{})
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

// validateJSONData validates the JSON string against the provided schema.
func validateJSONData(jsonSchema *gojsonschema.Schema, jsonString string) error {
	if err := ValidateJSON(jsonSchema, jsonString); err != nil {
		zap.L().Error("JSON validation failed", zap.Error(err))

		return err
	}

	zap.L().Info("JSON validation passed", zap.String("jsonString", jsonString))

	return nil
}

// Currently, it includes a stubbed-out comment for S3 logic.
func sendToS3(_ string) error {
	zap.L().Info("Stub: sendToS3 called.")

	return nil
}

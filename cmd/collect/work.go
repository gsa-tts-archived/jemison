package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/xeipuuv/gojsonschema"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

//nolint:revive
func (w *CollectWorker) Work(ctx context.Context, job *river.Job[common.CollectArgs]) error {
	zap.L().Info("Collect Worker processing job", zap.String("data-json", job.Args.Json))

	// Ensure schemas are initialized
	if err := ensureSchemasInitialized(); err != nil {
		return err
	}

	// Deserialize JSON data
	rawData, err := deserializeJSON(job.Args.Json)
	if err != nil {
		return err
	}

	// Check for `data.id` field and select schema
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
		zap.L().Error("failed to unmarshal JSON", zap.Error(err))
		return nil, err
	}
	return rawData, nil
}

func selectSchema(rawData map[string]interface{}) (*gojsonschema.Schema, error) {
	dataID, ok := rawData["id"].(string)
	if !ok || dataID == "" {
		return nil, errors.New("missing or invalid `data.id` field")
	}

	switch dataID {
	case "entree":
		return entreeSchema, nil
	case "fetch":
		return fetchSchema, nil
	default:
		return nil, errors.New("unsupported `data.id` value")
	}
}

func validateAndHandleBusinessLogic(schema *gojsonschema.Schema, args common.CollectArgs) error {
	if err := ValidateJSON(schema, args.Json); err != nil {
		zap.L().Error("JSON validation failed", zap.Error(err))
		return err
	}

	// Log validation success
	zap.L().Info("JSON validation passed",
		zap.String("json", args.Json))

	// Handle business logic
	return HandleBusinessLogic(args, args.Json)
}

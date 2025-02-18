package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/kv"
	"os"
	"path/filepath"
	"strings"
	"time"

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

func sendToS3(jsonString string) error {
	// Parse the source from the JSON for the naming convention
	source := "unknown"
	jsonData, err := deserializeJSON(jsonString)

	if err != nil {
		return fmt.Errorf("failed to parse source from JSON: %w", err)
	}

	if data, ok := jsonData["data"].(map[string]interface{}); ok {
		if src, ok := data["source"].(string); ok && src != "" {
			source = src
		}
	}

	// Generate the file name using the current epoch time and source
	epochTime := time.Now().Unix()
	fileName := fmt.Sprintf("%d_%s.json", epochTime, source)

	// Create a temporary file
	tempDir := os.TempDir()
	tempFilePath := filepath.Join(tempDir, fileName)

	// Write JSON content to the temporary file
	if err := os.WriteFile(tempFilePath, []byte(jsonString), 0644); err != nil {
		zap.L().Error("Failed to write JSON to a temporary file", zap.Error(err), zap.String("filePath", tempFilePath))
		return fmt.Errorf("failed to write JSON to temp file: %w", err)
	}

	zap.L().Info("Temporary file created successfully", zap.String("filePath", tempFilePath))

	// Get the S3 bucket for 'data'
	bucketName := "data"
	if !env.IsValidBucketName(bucketName) {
		errMsg := fmt.Sprintf("Invalid bucket name: %s", bucketName)
		zap.L().Error(errMsg)
		return fmt.Errorf(errMsg)
	}

	s3Client := kv.NewS3(bucketName)

	// Define an S3 key (path) for the file
	s3Key := filepath.Join(strings.TrimSuffix(source, "/"), fileName)

	// Upload the file to the S3 bucket
	mimeType := "application/json"
	if uploadErr := s3Client.FileToS3Path(s3Key, tempFilePath, mimeType); uploadErr != nil {
		zap.L().Error("Failed to upload file to S3", zap.Error(uploadErr), zap.String("s3Key", s3Key), zap.String("bucket", bucketName))
		return fmt.Errorf("failed to upload file to S3: %w", uploadErr)
	}

	zap.L().Info("Successfully uploaded file to S3", zap.String("s3Key", s3Key), zap.String("bucket", bucketName))

	// Clean up the temporary file
	if removeErr := os.Remove(tempFilePath); removeErr != nil {
		zap.L().Warn("Failed to delete temporary file", zap.Error(removeErr), zap.String("filePath", tempFilePath))
	} else {
		zap.L().Info("Temporary file deleted successfully", zap.String("filePath", tempFilePath))
	}

	return nil
}

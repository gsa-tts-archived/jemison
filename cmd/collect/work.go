package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/kv"
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

	// Pull in IsFull and hallpass
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
	// Parse the source from the JSON
	source, err := extractSourceFromJSON(jsonString)
	if err != nil {
		return err
	}

	// Create a temporary file for the JSON
	tempFilePath, err := createTempFile(jsonString, source)
	if err != nil {
		return err
	}
	defer cleanupTempFile(tempFilePath)

	// Validate S3 bucket name
	bucketName := "data"
	if !env.IsValidBucketName(bucketName) {
		zap.L().Error("Invalid bucket name", zap.String("bucketName", bucketName))

		return fmt.Errorf("invalid bucket name: %s", bucketName)
	}

	// Upload the file to S3
	s3Key := generateS3Key(source)
	if err := uploadFileToS3(bucketName, tempFilePath, s3Key); err != nil {
		return err
	}

	return nil
}

func extractSourceFromJSON(jsonString string) (string, error) {
	jsonData, err := deserializeJSON(jsonString)
	if err != nil {
		return "", fmt.Errorf("failed to parse source from JSON: %w", err)
	}

	if data, ok := jsonData["data"].(map[string]interface{}); ok {
		if src, ok := data["source"].(string); ok && src != "" {
			return src, nil // Return early if we find the `source`
		}
	}

	return "unknown", nil
}

func createTempFile(jsonString, source string) (string, error) {
	// Generate the file name using the current epoch time and source
	epochTime := time.Now().Unix()
	fileName := fmt.Sprintf("%d_%s.json", epochTime, source)
	tempDir := os.TempDir()
	tempFilePath := filepath.Join(tempDir, fileName)
	mode := 0o600

	if err := os.WriteFile(tempFilePath, []byte(jsonString), os.FileMode(mode)); err != nil {
		zap.L().Error("Failed to write JSON to a temporary file", zap.Error(err), zap.String("filePath", tempFilePath))

		return "", fmt.Errorf("failed to write JSON to temp file: %w", err)
	}

	zap.L().Info("Temporary file created successfully", zap.String("filePath", tempFilePath))

	return tempFilePath, nil
}

func cleanupTempFile(tempFilePath string) {
	if err := os.Remove(tempFilePath); err != nil {
		zap.L().Warn("Failed to delete temporary file", zap.Error(err), zap.String("filePath", tempFilePath))
	} else {
		zap.L().Info("Temporary file deleted successfully", zap.String("filePath", tempFilePath))
	}
}

func generateS3Key(source string) string {
	return filepath.Join(strings.TrimSuffix(source, "/"), fmt.Sprintf("%d_%s.json", time.Now().Unix(), source))
}

func uploadFileToS3(bucketName, tempFilePath, s3Key string) error {
	s3Client := kv.NewS3(bucketName)
	mimeType := "application/json"

	if err := s3Client.FileToS3Path(s3Key, tempFilePath, mimeType); err != nil {
		zap.L().Error("Failed to upload file to S3",
			zap.Error(err),
			zap.String("s3Key", s3Key),
			zap.String("bucket", bucketName),
		)

		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	zap.L().Info("Successfully uploaded file to S3", zap.String("s3Key", s3Key), zap.String("bucket", bucketName))

	return nil
}

package main

import (
	"context"
	"testing"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Test the validateJSONData function with valid input.
func TestValidateJSONData_ValidInput(t *testing.T) {
	// Initialize schemas
	err := InitializeSchemas()
	assert.NoError(t, err)

	// Valid JSON schema and instance
	jsonSchema := fetchSchema
	jsonString := `{
		"scheme": "https",
		"host": "example.gov",
		"path": "/api/resource",
		"data": { "id": "unique-fetch-id-5678", "source": "fetch", "payload": "example payload" }
	}`

	// Validate
	err = validateJSONData(jsonSchema, jsonString)
	assert.NoError(t, err, "Validation should pass for valid JSON input")
}

// Test the validateJSONData function with invalid input.
func TestValidateJSONData_InvalidInput(t *testing.T) {
	// Initialize schemas
	err := InitializeSchemas()
	assert.NoError(t, err)

	// Invalid JSON schema and instance (missing "data.source")
	jsonSchema := fetchSchema
	jsonString := `{
		"scheme": "https",
		"host": "example.gov",
		"path": "/api/resource",
		"data": { "id": "invalid-fetch-id", "payload": "example payload" }
	}`

	// Validation
	err = validateJSONData(jsonSchema, jsonString)
	assert.Error(t, err, "Validation should fail for invalid JSON input")
}

// Test for validating the Work function with a valid payload.
func TestWork_ValidFetchPayload(t *testing.T) {
	// Initialize schemas
	err := InitializeSchemas()
	assert.NoError(t, err)

	// Input
	args := common.CollectArgs{
		Scheme: "https",
		Host:   "example.gov",
		Path:   "/api/resource",
		JSON: `{
			"scheme": "https",
			"host": "example.gov",
			"path": "/api/resource",
			"api-key": "api-key-123",
			"data": { "id": "unique-fetch-id-5678", "source": "fetch", "payload": "example payload" }
		}`,
	}

	// Create a mock job
	job := &river.Job[common.CollectArgs]{Args: args}

	// Mock logger
	loggerConfig := zap.NewDevelopmentConfig()
	logger, _ := loggerConfig.Build()
	zap.ReplaceGlobals(logger)

	worker := &CollectWorker{}
	err = worker.Work(context.Background(), job)

	// Validate the worker's success
	assert.NoError(t, err, "Worker should process valid fetch payload successfully")
}

// Test for invalid schema in payload.
func TestWork_InvalidSchemaPayload(t *testing.T) {
	err := InitializeSchemas()
	assert.NoError(t, err)

	// Invalid payload does not match schema (missing required key: data.source)
	args := common.CollectArgs{
		JSON: `{"scheme":"https","host":"example.gov","path":"/api/resource","data":{"id":"fetch-id","payload":"invalid format"}}`, //nolint:lll
	}

	job := &river.Job[common.CollectArgs]{Args: args}

	// Mock logger
	loggerConfig := zap.NewDevelopmentConfig()
	logger, _ := loggerConfig.Build()
	zap.ReplaceGlobals(logger)

	worker := &CollectWorker{}
	err = worker.Work(context.Background(), job)

	assert.Error(t, err, "Worker should return an error for invalid payload schema")
}

// Test for unknown source value.
func TestWork_UnknownSource(t *testing.T) {
	err := InitializeSchemas()
	assert.NoError(t, err)

	// Unknown source field in data
	args := common.CollectArgs{
		JSON: `{"scheme":"https","host":"example.gov","path":"/api/resource","data":{"id":"unknown-id","source":"unknown","payload":"test"}}`, //nolint:lll
	}

	job := &river.Job[common.CollectArgs]{Args: args}

	// Mock logger
	loggerConfig := zap.NewDevelopmentConfig()
	logger, _ := loggerConfig.Build()
	zap.ReplaceGlobals(logger)

	worker := &CollectWorker{}
	err = worker.Work(context.Background(), job)

	assert.Error(t, err, "Worker should encounter error for unknown `source` value")
}

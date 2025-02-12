package main

import (
	"context"
	"testing"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type MockWorker struct{}

func (w *MockWorker) Work(ctx context.Context, job *river.Job[common.CollectArgs]) error {
	return (&CollectWorker{}).Work(ctx, job)
}

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

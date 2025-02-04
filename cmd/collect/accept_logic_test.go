package main

import (
	"testing"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTransformArgumentsToJSON(t *testing.T) {
	// Sample input
	args := common.CollectArgs{
		Scheme: "https",
		Host:   "www.example.com",
		Path:   "/test",
	}

	jsonString, err := TransformArgumentsToJSON(args)

	// Validate
	assert.NoError(t, err, "TransformArgumentsToJSON should not return an error")
	assert.JSONEq(t, `{"Scheme":"https","Host":"www.example.com","Path":"/test"}`, jsonString, "JSON output is incorrect")
}

func TestHandleBusinessLogic(t *testing.T) {
	// Mock zap logger
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	logger, _ := loggerConfig.Build()

	defer func() {
		if err := logger.Sync(); err != nil {
			t.Logf("failed to sync logger: %v", err)
		}
	}()

	zap.ReplaceGlobals(logger)

	// Sample input
	args := common.CollectArgs{
		Scheme: "http",
		Host:   "example.org",
		Path:   "/example",
	}

	err := HandleBusinessLogic(args)

	// Validate
	assert.NoError(t, err, "HandleBusinessLogic should not return an error")
}

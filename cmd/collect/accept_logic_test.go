package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestValidateJSON(t *testing.T) {
	// Set up valid and invalid JSON examples
	validJSON := `{"scheme":"https","host":"example.gov","path":"/api/resource","api-key":"key123","data":{"id":"unique-fetch-id-5678","source":"fetch","payload":"some payload"}}` //nolint:lll
	invalidJSON := `{"scheme":"https","host":"example.gov","path":"/api/resource","data":{"id":"123"}}`

	// Fetch schema should be already initialized for this test
	err := InitializeSchemas()
	assert.NoError(t, err, "Schema initialization should not fail")

	// Test valid JSON
	err = ValidateJSON(fetchSchema, validJSON)
	assert.NoError(t, err, "Valid JSON should pass validation")

	// Test invalid JSON
	err = ValidateJSON(fetchSchema, invalidJSON)
	assert.Error(t, err, "Invalid JSON should fail validation")
}

func TestInitializeSchemas(t *testing.T) {
	// First initialization
	err := InitializeSchemas()
	assert.NoError(t, err, "First schema initialization should succeed")

	// Subsequent initialization call (idempotency check)
	err = InitializeSchemas()
	assert.NoError(t, err, "Subsequent schema initialization should succeed")
}

func TestSelectSchema(t *testing.T) {
	// Mock zap logger
	loggerConfig := zap.NewDevelopmentConfig()
	logger, _ := loggerConfig.Build()
	zap.ReplaceGlobals(logger)

	// Valid entree JSON
	entreeData := map[string]interface{}{
		"data": map[string]interface{}{
			"source": "entree",
		},
	}

	// Valid fetch JSON
	fetchData := map[string]interface{}{
		"data": map[string]interface{}{
			"source": "fetch",
		},
	}

	// Invalid schema JSON
	invalidData := map[string]interface{}{
		"data": map[string]interface{}{
			"source": "unknown",
		},
	}

	// Check schema selection
	schema, err := selectSchema(entreeData)
	assert.NoError(t, err)
	assert.Equal(t, entreeSchema, schema)

	schema, err = selectSchema(fetchData)
	assert.NoError(t, err)
	assert.Equal(t, fetchSchema, schema)

	_, err = selectSchema(invalidData)
	assert.Error(t, err, "Unknown source should result in error")
}

package common

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

// MarshalMapToJSON marshals a map into a JSON string and logs any errors.
func MarshalMapToJSON(data map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to marshal map to JSON: %w", err)
		zap.L().Error(wrappedErr.Error(), zap.Error(err))

		return "", wrappedErr
	}

	return string(jsonData), nil
}

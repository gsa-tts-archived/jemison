package main

import (
	"encoding/json"
	"fmt"

	"github.com/GSA-TTS/jemison/internal/common"
	"go.uber.org/zap"
)

func TransformArgumentsToJSON(args common.CollectArgs) (string, error) {
	jsonData, err := json.Marshal(args)
	if err != nil {
		return "", fmt.Errorf("failed to marshal CollectArgs to JSON: %w", err)
	}

	return string(jsonData), nil
}

func HandleBusinessLogic(args common.CollectArgs) error {
	jsonString, err := TransformArgumentsToJSON(args)
	if err != nil {
		zap.L().Error("failed to transform arguments to JSON", zap.Error(err))

		return err
	}

	zap.L().Info("Collect service captured crawl record", zap.String("json", jsonString))

	return nil
}

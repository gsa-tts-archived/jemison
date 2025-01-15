package main

import (
	"context"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/riverqueue/river"
)

//nolint:revive
func (w *ResultsApiWorker) Work(ctx context.Context, job *river.Job[common.ResultsApiArgs]) error {
	return nil
}

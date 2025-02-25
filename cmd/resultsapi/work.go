package main

import (
	"context"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/riverqueue/river"
)

//nolint:revive
func (w *ResultsAPIWorker) Work(ctx context.Context, job *river.Job[common.ResultsAPIArgs]) error {
	return nil
}

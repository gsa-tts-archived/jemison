package main

import (
	"context"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/riverqueue/river"
)

func (w *ServeWorker) Work(ctx context.Context, job *river.Job[common.ServeArgs]) error {
	return nil
}

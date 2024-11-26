package main

import (
	"context"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/riverqueue/river"
)

func (w *BalanceWorker) Work(ctx context.Context, job *river.Job[common.BalanceArgs]) error {

	return nil
}

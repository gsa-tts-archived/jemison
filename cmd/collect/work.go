package main

import (
	"context"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

//func (w *CollectWorker) Work(ctx context.Context, job *river.Job[common.CollectArgs]) error {
//	return nil
//}

//nolint:revive
func (w *CollectWorker) Work(ctx context.Context, job *river.Job[common.CollectArgs]) error {
	zap.L().Info("Collect has access to crawl record",
		zap.String("scheme", job.Args.Scheme),
		zap.String("host", job.Args.Host),
		zap.String("path", job.Args.Path))

	if err := HandleBusinessLogic(job.Args); err != nil {
		zap.L().Error("failed to handle business logic", zap.Error(err))

		return err
	}

	return nil
}

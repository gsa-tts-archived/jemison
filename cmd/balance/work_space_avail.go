package main

import (
	"context"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

// The SpaceAvail worker handles the case where a serve service is announcing
// how much space it has available. We use this to build a map of who has what
// free space.

func (w *SpaceAvailWorker) Work(ctx context.Context, job *river.Job[common.SpaceAvailArgs]) error {
	err := TheSpaceMap.Add(job.Args.Filename, job.Args.Size)
	if err != nil {
		zap.L().Error("do not have space for this database",
			zap.String("filename", job.Args.Filename),
			zap.Int64("size", job.Args.Size))

		// Requeue to someone else?
		return nil
	}
	return nil
}

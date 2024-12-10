package main

import (
	"context"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/kv"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

func (w *ServeWorker) Work(ctx context.Context, job *river.Job[common.ServeArgs]) error {

	s, _ := env.Env.GetUserService(ThisServiceName)

	databases_file_path := s.GetParamString("database_files_path")

	sqlite_filename := sqlite.SqliteFilename(job.Args.Filename)

	zap.L().Debug("received sqlite filename",
		zap.String("filename", sqlite_filename))

	// Writes to the local filesystem.
	destination := databases_file_path + "/" + sqlite_filename

	zap.L().Info("copying file to host",
		zap.String("destination", destination),
		zap.String("sqlite_filename", sqlite_filename),
	)

	// Downloads content to the destination
	s3 := kv.NewS3(ThisServiceName)
	s3.S3PathToFile(sqlite_filename, destination)
	return nil
}

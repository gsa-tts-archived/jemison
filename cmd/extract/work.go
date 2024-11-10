package main

import (
	"context"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/GSA-TTS/jemison/internal/kv"
	"github.com/GSA-TTS/jemison/internal/util"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

func extract(obj *kv.S3JSON) {
	mime_type := obj.GetString("content-type")
	s, _ := env.Env.GetUserService(ThisServiceName)

	switch mime_type {
	case "text/html":
		if s.GetParamBool("extract_html") {
			extractHtml(obj)
		}
	case "application/pdf":
		if s.GetParamBool("extract_pdf") {
			extractPdf(obj)
		}
	}
}

func (w *ExtractWorker) Work(ctx context.Context, job *river.Job[common.ExtractArgs]) error {

	zap.L().Debug("extracting",
		zap.String("host", job.Args.Host),
		zap.String("path", job.Args.Path))

	s3json := kv.NewEmptyS3JSON("fetch", util.ToScheme(job.Args.Scheme), job.Args.Host, job.Args.Path)
	s3json.Load()

	extract(s3json)

	zap.L().Debug("extraction finished")
	return nil

}

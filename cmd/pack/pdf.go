package main

import (
	"github.com/GSA-TTS/jemison/internal/kv"
	"github.com/GSA-TTS/jemison/internal/sqlite"
	"github.com/GSA-TTS/jemison/internal/sqlite/schemas"
	"go.uber.org/zap"
)

func packPdf(pt *sqlite.PackTable, s3json *kv.S3JSON) {
	// We have more fields than before.
	path_id, err := pt.Queries.InsertPath(pt.Context, s3json.GetString("path"))
	if err != nil {
		zap.L().Error("could not insert path when packing PDF", zap.String("path", s3json.GetString("path")))
	}
	// Insert the title
	_, err = pt.Queries.InsertTitle(pt.Context, schemas.InsertTitleParams{
		PathID: path_id,
		Title:  s3json.GetString("title") + " (PDF page " + s3json.GetString("pdf_page_number") + ")",
	})
	if err != nil {
		zap.L().Error("could not insert title when packing PDF",
			zap.String("title", s3json.GetString("title")))
	}
	// Insert the content
	_, err = pt.Queries.InsertBody(pt.Context, schemas.InsertBodyParams{
		PathID: path_id,
		Body:   s3json.GetString("content"),
	})
	if err != nil {
		zap.L().Error("could not insert body when packing PDF",
			zap.String("title", s3json.GetString("title")))
	}
}

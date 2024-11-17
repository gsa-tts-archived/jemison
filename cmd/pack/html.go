package main

import (
	"github.com/GSA-TTS/jemison/internal/kv"
	"github.com/GSA-TTS/jemison/internal/sqlite"
	"github.com/GSA-TTS/jemison/internal/sqlite/schemas"
	"go.uber.org/zap"
)

func packHtml(pt *sqlite.PackTable, s3json *kv.S3JSON) {
	// We have more fields than before.
	path_id, err := pt.Queries.InsertPath(pt.Context, s3json.GetString("path"))
	if err != nil {
		zap.L().Error("could not insert path when packing",
			zap.String("_key", s3json.GetString("_key")),
			zap.String("err", err.Error()),
			zap.String("path", s3json.GetString("path")),
		)
	}

	// Insert the title
	_, err = pt.Queries.InsertTitle(pt.Context, schemas.InsertTitleParams{
		PathID: path_id,
		Title:  s3json.GetString("title"),
	})
	if err != nil {
		zap.L().Error("could not insert title when packing",
			zap.String("_key", s3json.GetString("_key")),
			zap.String("err", err.Error()),
			zap.String("title", s3json.GetString("title")))
	}

	// Insert the content
	body := s3json.GetString("body")
	_, err = pt.Queries.InsertBody(pt.Context, schemas.InsertBodyParams{
		PathID: path_id,
		Body:   body,
	})
	if err != nil {
		zap.L().Error("could not insert body when packing",
			zap.String("_key", s3json.GetString("_key")),
			zap.String("err", err.Error()),
			zap.String("sqlite_filename", pt.Filename),
			zap.String("title", s3json.GetString("title")),
			zap.String("body", body[:min(len(body), 30)]+"..."),
		)
	}
}

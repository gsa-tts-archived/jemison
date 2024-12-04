package main

import (
	"context"

	"github.com/GSA-TTS/jemison/internal/search_db/search_db"

	"github.com/GSA-TTS/jemison/internal/kv"
	"github.com/GSA-TTS/jemison/internal/sqlite"
	"github.com/GSA-TTS/jemison/internal/sqlite/schemas"
	"go.uber.org/zap"
)

func packHtml(pt *sqlite.PackTable, s3json *kv.S3JSON) {
	// We have more fields than before.
	path_id, err := pt.Queries.InsertPath(context.Background(), s3json.GetString("path"))
	if err != nil {
		zap.L().Error("could not insert path when packing",
			zap.String("_key", s3json.GetString("_key")),
			zap.String("err", err.Error()),
			zap.String("path", s3json.GetString("path")),
		)
	}

	// Insert the title
	_, err = pt.Queries.InsertTitle(context.Background(), schemas.InsertTitleParams{
		PathID: path_id,
		Title:  s3json.GetString("title"),
	})
	if err != nil {
		zap.L().Error("could not insert title when packing",
			zap.String("_key", s3json.GetString("_key")),
			zap.String("err", err.Error()),
			zap.String("title", s3json.GetString("title")))
	}

	// Insert the headers
	_, err = pt.Queries.InsertHeader(context.Background(), schemas.InsertHeaderParams{
		PathID: path_id,
		// FIXME: There's a level hiding in here...
		Header: s3json.GetString("header"),
	})
	if err != nil {
		zap.L().Error("could not insert header when packing",
			zap.String("_key", s3json.GetString("_key")),
			zap.String("err", err.Error()),
			zap.String("header", s3json.GetString("header")))
	}

	// Insert the content
	body := s3json.GetString("body")
	_, err = pt.Queries.InsertBody(context.Background(), schemas.InsertBodyParams{
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

func packPgHtml(sdb *search_db.SearchDB, s3json *kv.S3JSON) {
	// We have more fields than before.
	path_id, err := sdb.Queries.InsertPath(context.Background(), s3json.GetString("path"))
	if err != nil {
		zap.L().Error("could not insert path when packing",
			zap.String("_key", s3json.GetString("_key")),
			zap.String("err", err.Error()),
			zap.String("path", s3json.GetString("path")),
		)
	}

	// Insert the title
	_, err = sdb.Queries.InsertTitle(context.Background(), search_db.InsertTitleParams{
		PathID: int32(path_id),
		Title:  s3json.GetString("title"),
	})
	if err != nil {
		zap.L().Error("could not insert title when packing",
			zap.String("_key", s3json.GetString("_key")),
			zap.String("err", err.Error()),
			zap.String("title", s3json.GetString("title")))
	}

	// Insert the headers
	_, err = sdb.Queries.InsertHeader(context.Background(), search_db.InsertHeaderParams{
		PathID: int32(path_id),
		// FIXME: There's a level hiding in here...
		Header: s3json.GetString("header"),
	})
	if err != nil {
		zap.L().Error("could not insert header when packing",
			zap.String("_key", s3json.GetString("_key")),
			zap.String("err", err.Error()),
			zap.String("header", s3json.GetString("header")))
	}

	// Insert the content
	body := s3json.GetString("body")
	_, err = sdb.Queries.InsertBody(context.Background(), search_db.InsertBodyParams{
		PathID: int32(path_id),
		Body:   body,
	})
	if err != nil {
		zap.L().Error("could not insert body when packing",
			zap.String("_key", s3json.GetString("_key")),
			zap.String("err", err.Error()),
			zap.String("title", s3json.GetString("title")),
			zap.String("body", body[:min(len(body), 30)]+"..."),
		)
	}
}

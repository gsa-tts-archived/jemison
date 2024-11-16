package queueing

import (
	"fmt"

	"github.com/GSA-TTS/jemison/internal/common"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"go.uber.org/zap"
)

func InsertServe(sqlite_s3_key string) {
	_, pool, _ := common.CommonQueueInit()
	ctx, tx := common.CtxTx(pool)
	defer pool.Close()
	defer tx.Commit(ctx)

	client, err := river.NewClient(riverpgxv5.New(pool), &river.Config{})
	if err != nil {
		zap.L().Error("could not create river client",
			zap.String("error", err.Error()))
	}
	client.InsertTx(ctx, tx, common.ServeArgs{
		Filename: sqlite_s3_key,
	}, &river.InsertOpts{Queue: "serve"})
	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		zap.L().Panic("cannot commit insert tx",
			zap.String("filename", sqlite_s3_key))
	}
}

func InsertFetch(scheme string, host string, path string) {
	_, pool, _ := common.CommonQueueInit()
	ctx, tx := common.CtxTx(pool)
	defer pool.Close()
	defer tx.Commit(ctx)

	client, err := river.NewClient(riverpgxv5.New(pool), &river.Config{})
	if err != nil {
		zap.L().Error("could not create river client",
			zap.String("error", err.Error()))
	}
	client.InsertTx(ctx, tx, common.FetchArgs{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}, &river.InsertOpts{Queue: "fetch"})
	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		zap.L().Panic("cannot commit insert tx",
			zap.String("host", host),
			zap.String("path", path))
	}
}

func InsertValidate(queue string, job river.JobArgs) {
	_, pool, _ := common.CommonQueueInit()
	ctx, tx := common.CtxTx(pool)
	defer pool.Close()
	defer tx.Commit(ctx)

	client, err := river.NewClient(riverpgxv5.New(pool), &river.Config{})
	if err != nil {
		zap.L().Error("could not create river client",
			zap.String("error", err.Error()))
	}
	client.InsertTx(ctx, tx, job, &river.InsertOpts{Queue: queue})
	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		zap.L().Panic("cannot commit insert tx",
			zap.String("job", fmt.Sprintln(job)))
	}
}

func InsertExtract(scheme string, host string, path string) {
	_, pool, _ := common.CommonQueueInit()
	ctx, tx := common.CtxTx(pool)
	defer pool.Close()
	defer tx.Commit(ctx)

	client, err := river.NewClient(riverpgxv5.New(pool), &river.Config{})
	if err != nil {
		zap.L().Error("could not create river client",
			zap.String("error", err.Error()))
	}
	client.InsertTx(ctx, tx, common.ExtractArgs{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}, &river.InsertOpts{Queue: "extract"})
	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		zap.L().Panic("cannot commit insert tx",
			zap.String("host", host),
			zap.String("path", path))
	}
}

func InsertWalk(scheme string, host string, path string) {
	_, pool, _ := common.CommonQueueInit()
	ctx, tx := common.CtxTx(pool)
	defer pool.Close()
	defer tx.Commit(ctx)

	client, err := river.NewClient(riverpgxv5.New(pool), &river.Config{})
	if err != nil {
		zap.L().Error("could not create river client",
			zap.String("error", err.Error()))
	}
	client.InsertTx(ctx, tx, common.WalkArgs{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}, &river.InsertOpts{Queue: "walk"})
	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		zap.L().Panic("cannot commit insert tx",
			zap.String("host", host),
			zap.String("path", path))
	}
}

func InsertEntree(scheme string, host string, path string, isFullCrawl bool, isHallPass bool) {
	_, pool, _ := common.CommonQueueInit()
	ctx, tx := common.CtxTx(pool)
	defer pool.Close()
	defer tx.Commit(ctx)

	client, err := river.NewClient(riverpgxv5.New(pool), &river.Config{})
	if err != nil {
		zap.L().Error("could not create river client",
			zap.String("error", err.Error()))
	}
	client.InsertTx(ctx, tx, common.EntreeArgs{
		Scheme:    scheme,
		Host:      host,
		Path:      path,
		FullCrawl: isFullCrawl,
		HallPass:  isHallPass,
	}, &river.InsertOpts{Queue: "entree"})
	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		zap.L().Panic("cannot commit insert tx",
			zap.String("host", host),
			zap.String("path", path),
			zap.Bool("hallpass", isHallPass),
		)
	}
}

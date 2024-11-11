package queueing

import (
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

package queueing

import (
	"context"
	"time"

	"github.com/GSA-TTS/jemison/internal/common"
	"go.uber.org/zap"
)

func ClearCompletedPeriodically() {
	_, pool, _ := common.CommonQueueInit()
	defer pool.Close()

	ticker := time.NewTicker(3 * time.Minute)
	for {
		<-ticker.C
		zap.L().Warn("clearing completed queue")
		ctx := context.Background()
		pool.Exec(ctx, "DELETE FROM river_job WHERE state='completed'")
	}
}

package queueing

import (
	"github.com/riverqueue/river"
)

type GenericRequest struct {
	Key       string `json:"key"`
	QueueName string
}

func NewGenericRequest() GenericRequest {
	return GenericRequest{}
}

func (g GenericRequest) Kind() string {
	return g.QueueName
}

func (g GenericRequest) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: g.Kind(),
	}
}

// func InsertFetch(fetch_args common.FetchArgs) {

// 	zap.L().Debug("Inserting validate job")
// 	ctx2, tx2 := common.CtxTx(walkPool)
// 	defer tx2.Rollback(ctx)
// 	walkClient.InsertTx(ctx2, tx2, common.ValidateFetchArgs{
// 		Fetch: common.FetchArgs{
// 			Scheme: job.Args.Scheme,
// 			Host:   job.Args.Host,
// 			Path:   job.Args.Path,
// 		},
// 	}, &river.InsertOpts{Queue: "validate_fetch"})
// 	if err := tx2.Commit(ctx2); err != nil {
// 		zap.L().Panic("cannot commit validate tx",
// 			zap.String("error", err.Error()))
// 	}
// }

// // To be flexible, and use this across the app...
// // this process takes a channel of maps. It looks at the map, and
// // then routes it to the right River work queue.
// func Enqueue(generic_job chan<- kv.JSON) {

// 	// Set up the River client

// 	for {
// 		// Switch on the job kind

// 	}
// }

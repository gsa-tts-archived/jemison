package main

import (
	"fmt"
	"runtime"

	"github.com/GSA-TTS/jemison/internal/common"
	kv "github.com/GSA-TTS/jemison/internal/kv"
	"github.com/GSA-TTS/jemison/internal/util"
	"github.com/johbar/go-poppler"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

func extractPdf(obj *kv.S3JSON) {
	rawFilename := obj.GetString("raw")
	s3 := kv.NewS3(ThisServiceName)
	s3.S3ToFile(obj.Key, rawFilename)

	doc, err := poppler.Open(rawFilename)

	if err != nil {
		fmt.Println("Failed to convert body to Document")
	} else {
		for page_no := 0; page_no < doc.GetNPages(); page_no++ {

			page_number_anchor := fmt.Sprintf("#page=%d", page_no+1)
			copied_key := obj.Key.Copy()
			copied_key.Path = copied_key.Path + page_number_anchor
			copied_key.Extension = util.JSON

			page := doc.GetPage(page_no)
			obj.Set("content", util.RemoveStopwords(page.Text()))
			obj.Set("path", copied_key.Path)
			obj.Set("pdf_page_number", fmt.Sprintf("%d", page_no+1))

			new_obj := kv.NewFromBytes(
				ThisServiceName,
				obj.Key.Scheme,
				obj.Key.Host,
				obj.Key.Path,
				obj.GetJSON(),
			)
			new_obj.Save()
			page.Close()
			// e.Stats.Increment("page_count")

			// Enqueue next steps
			ctx, tx := common.CtxTx(packPool)
			defer tx.Rollback(ctx)

			packClient.InsertTx(ctx, tx, common.PackArgs{
				Scheme: obj.Key.Scheme.String(),
				Host:   obj.Key.Host,
				Path:   obj.Key.Path,
			}, &river.InsertOpts{Queue: "pack"})
			if err := tx.Commit(ctx); err != nil {
				zap.L().Panic("cannot commit insert tx",
					zap.String("key", obj.Key.Render()))
			}

			// https://weaviate.io/blog/gomemlimit-a-game-changer-for-high-memory-applications
			// https://stackoverflow.com/questions/38972003/how-to-stop-the-golang-gc-and-trigger-it-manually
			runtime.GC()
		}
	}

	//e.Stats.Increment("document_count")

	doc.Close()
}

package main

import (
	"os"
	"strings"

	"github.com/GSA-TTS/jemison/internal/common"
	kv "github.com/GSA-TTS/jemison/internal/kv"
	"github.com/GSA-TTS/jemison/internal/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

func scrape_sel(sel *goquery.Selection) string {
	content := ""
	txt := sel.Text()
	if len(txt) > 0 {
		repl := strings.ToLower(txt)
		// FIXME: This should be part of a standalone text processing
		// module. Perhaps we don't always want to do this. For example,
		// in the case of multi-lingual scraping.
		repl = util.RemoveStopwords(repl)
		repl += " "
		if len(repl) > 2 {
			content += repl
		}
	}
	return content
}

func extractHtml(obj *kv.S3JSON) {
	// rawFilename := obj.GetString("raw")
	rawFilename := uuid.NewString()
	// The file is not in this service... it's in the `fetch` bucket.`
	s3 := kv.NewS3("fetch")

	raw_key := obj.Key.Copy()
	raw_key.Extension = util.Raw
	s3.S3ToFile(raw_key, rawFilename)
	rawFile, err := os.Open(rawFilename)
	if err != nil {
		zap.L().Error("cannot open tempfile", zap.String("filename", rawFilename))
	}
	defer func() {
		rawFile.Close()
		os.Remove(rawFilename)
	}()

	doc, err := goquery.NewDocumentFromReader(rawFile)
	if err != nil {
		zap.L().Fatal("cannot create new doc from raw file")
	}

	// fi, err := os.Stat(rawFilename)
	// if err != nil {
	// 	zap.L().Fatal(err.Error())
	// }
	// // get the size
	// size := fi.Size()
	// zap.L().Debug("size", zap.Int64("size", size))

	content := ""
	for _, elem := range []string{
		"p",
		"li",
		"td",
		"div",
		"span",
		"a",
		"small",
		"b",
		"bold",
		"em",
		"i",
		"title",
		"h1",
		"h2",
		"h3",
		"h4",
		"h5",
		"h6",
		"h7",
		"h8",
	} {
		// zap.L().Debug("looking for", zap.String("elem", elem))
		doc.Find(elem).Each(func(ndx int, sel *goquery.Selection) {
			// zap.L().Debug("element", zap.String("sel", scrape_sel(sel)))
			content += scrape_sel(sel)
		})
	}

	// Store everything
	copied_key := obj.Key.Copy()
	copied_key.Extension = util.JSON
	// This is because we were holding an object from the "fetch" bucket.
	new_obj := kv.NewFromBytes(
		ThisServiceName,
		obj.Key.Scheme,
		obj.Key.Host,
		obj.Key.Path,
		obj.GetJSON())
	// zap.L().Debug("content check", zap.String("content", util.RemoveStopwords(content)))
	new_obj.Set("content", util.RemoveStopwords(content))
	new_obj.Save()

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

}

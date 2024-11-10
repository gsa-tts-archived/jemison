package main

import (
	"log"
	"os"
	"strings"

	"github.com/GSA-TTS/jemison/internal/common"
	kv "github.com/GSA-TTS/jemison/internal/kv"
	"github.com/GSA-TTS/jemison/internal/util"
	"github.com/PuerkitoBio/goquery"
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
	rawFilename := obj.GetString("raw")
	s3 := kv.NewS3(ThisServiceName)
	s3.S3ToFile(obj.Key, rawFilename)
	rawFile, err := os.Open(rawFilename)
	if err != nil {
		zap.L().Error("cannot open tempfile", zap.String("filename", rawFilename))
	}
	defer rawFile.Close()

	content := ""

	doc, err := goquery.NewDocumentFromReader(rawFile)
	if err != nil {
		log.Println("HTML cannot create new document")
		log.Fatal(err)
	}

	for _, elem := range [...]string{
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
		doc.Find(elem).Each(func(ndx int, sel *goquery.Selection) {
			content += scrape_sel(sel)
		})
	}

	// Store everything
	copied_key := obj.Key.Copy()
	copied_key.Extension = util.JSON
	obj.Set("content", util.RemoveStopwords(content))
	// This is because we were holding an object from the "fetch" bucket.
	new_obj := kv.NewFromBytes(
		ThisServiceName,
		obj.Key.Scheme,
		obj.Key.Host,
		obj.Key.Path,
		obj.GetJSON())
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

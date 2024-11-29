package main

import (
	"bytes"
	"context"
	"errors"
	"net/url"
	"os"
	"strings"

	common "github.com/GSA-TTS/jemison/internal/common"
	"github.com/GSA-TTS/jemison/internal/kv"
	"github.com/GSA-TTS/jemison/internal/queueing"
	"github.com/GSA-TTS/jemison/internal/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

// //////////////////////////////////////
// go_for_a_walk
func go_for_a_walk(s3json *kv.S3JSON) {
	cleaned_mime_type := util.CleanMimeType(s3json.GetString("content-type"))
	switch cleaned_mime_type {
	case "text/html":
		walk_html(s3json)
	case "application/pdf":
		// log.Println("PDFs do not walk")
	}
}

// //////////////////////////////////////
// extract_links
func extract_links(s3json *kv.S3JSON) []*url.URL {

	// Return a unique set
	link_set := make(map[string]bool)
	// Remove all trailing slashes.
	links := make([]*url.URL, 0)

	raw := s3json.GetString("raw")
	// This is a key to a file.
	tempFilename := uuid.NewString()

	s3 := kv.NewS3("fetch")
	err := s3.S3PathToFile(raw, tempFilename)
	if err != nil {
		zap.L().Error("could not copy s3 to local file",
			zap.String("tempFilename", tempFilename),
			zap.String("raw", raw))
		return links
	}

	tFile, err := os.Open(tempFilename)
	if err != nil {
		zap.L().Error("cannot open temporary file", zap.String("filename", tempFilename))
		return links
	}
	defer func() { tFile.Close(); os.Remove(tempFilename) }()

	fi, err := os.Stat(tempFilename)
	if err != nil {
		zap.L().Error(err.Error())
		return links
	}
	// get the size
	size := fi.Size()
	zap.L().Debug("tempFilename size", zap.Int64("size", size))

	doc, err := goquery.NewDocumentFromReader(tFile)

	if err != nil {
		zap.L().Error("WALK cannot convert to document", zap.String("key", s3json.Key.Render()))
		return links
	}

	// Pages that are javascript-only, or require javascript to display will
	// have no links.
	// zap.L().Debug("doc", zap.String("all", fmt.Sprintln(doc.Text())))

	doc.Find("a[href]").Each(func(ndx int, sel *goquery.Selection) {
		link, exists := sel.Attr("href")
		// zap.L().Debug("found link", zap.String("link", link), zap.Bool("exists", exists))

		if exists {
			link_to_crawl, err := is_crawlable(s3json, link)
			if err != nil {
				zap.L().Debug("error checking crawlability",
					zap.String("url", link),
					zap.String("error", err.Error()))
			} else {
				if _, ok := expirable_cache.Get(link_to_crawl); ok {
					// PASS ON LOGGING IF IT IS A CACHE HIT
				} else {
					// CRAWL BOTH HTTPS AND HTTP?
					if strings.HasPrefix(link_to_crawl, "http") {
						zap.L().Debug("link to crawl", zap.String("url", link_to_crawl))
						expirable_cache.Set(link_to_crawl, 0, 0)
						link_set[link_to_crawl] = true
					}
				}
			}
		}
	})

	for link := range link_set {
		link = trimSuffix(link, "/")
		u, err := url.Parse(link)
		if err != nil {
			zap.L().Warn("WALK ExtractLinks did a bad with",
				zap.String("badLink", link))
		} else {
			links = append(links, u)
		}
	}

	//log.Println("EXTRACTED", links)
	return links
}

// //////////////////////////////////////
// walk_html
func walk_html(s3json *kv.S3JSON) {
	links := extract_links(s3json)
	zap.L().Debug("walk considering links",
		zap.Int("count", len(links)))
	for _, link := range links {
		// The links come back canonicalized against the host. So,
		// all the fields should be present.
		ChQSHP <- queueing.QSHP{
			Queue:      "entree",
			Scheme:     link.Scheme,
			Host:       link.Host,
			Path:       link.Path,
			IsFull:     false,
			IsHallPass: false,
		}
	}
}

func is_crawlable(s3json *kv.S3JSON, link string) (string, error) {
	base := url.URL{
		Scheme: s3json.GetString("scheme"),
		Host:   s3json.GetString("host"),
		Path:   s3json.GetString("path"),
	}

	// zap.L().Debug("considering the url",
	// 	zap.String("url", link))

	// Is the URL at least length 1?
	if len(link) < 1 {
		return "", errors.New("crawler: URL is too short to crawl")
	}

	// Is it a relative URL? Then it is OK.
	if strings.HasPrefix(link, "/") || !strings.HasPrefix(link, "http") {
		u, _ := url.Parse(link)
		u = base.ResolveReference(u)
		zap.L().Debug("looks good to crawl",
			zap.String("base", base.String()),
			zap.String("resolved", u.String()))
		return u.String(), nil
	}

	lu, err := url.Parse(link)
	if err != nil {
		return "", errors.New("crawler: link does not parse")
	}

	// Does it end in .gov?
	// if bytes.HasSuffix([]byte(lu.Host), []byte("gov")) {
	// 	return "", errors.New("crawler: URL does not end in .gov")
	// }

	pieces := strings.Split(base.Host, ".")
	if len(pieces) < 2 {
		return "", errors.New("crawler: link host has too few pieces")
	} else {
		tld := pieces[len(pieces)-2] + "." + pieces[len(pieces)-1]

		// Does the link contain our TLD?
		if !strings.Contains(lu.Host, tld) {
			return "", errors.New("crawler: link does not contain the TLD")
		}
	}

	// FIXME: There seem to be whitespace URLs coming through. I don't know why.
	// This could be revisited, as it is expensive.
	// Do we still have garbage?
	if !bytes.HasPrefix([]byte(lu.String()), []byte("https")) ||
		!bytes.HasPrefix([]byte(lu.String()), []byte("http")) {
		return "", errors.New("crawler: link does not start with http(s)")
	}
	// Is it pure whitespace?
	if len(strings.Replace(lu.String(), " ", "", -1)) < 5 {
		return "", errors.New("crawler: link too short")
	}
	return lu.String(), nil
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
		return s
	} else {
		return s
	}
}

func (w *WalkWorker) Work(ctx context.Context, job *river.Job[common.WalkArgs]) error {
	if job.Attempt > 2 {
		zap.L().Warn("walking zombie; dropping",
			zap.String("host", job.Args.Host), zap.String("path", job.Args.Path))
		return nil
	}

	zap.L().Debug("walk job",
		zap.String("scheme", job.Args.Scheme),
		zap.String("host", job.Args.Host),
		zap.String("path", job.Args.Path),
	)

	s3json := kv.NewEmptyS3JSON("fetch",
		util.ToScheme(job.Args.Scheme),
		job.Args.Host,
		job.Args.Path)
	err := s3json.Load()
	if err != nil {
		// Don't do anything if we can't load the S3.
		// Quietly skip
		// FIXME: log this in the future
		return nil
	}
	// If we're here, we already fetched the content.
	// So, add ourselves to the cache. Don't re-crawl ourselves
	// FIXME: figure out if the scheme ends up in the JSON
	expirable_cache.Set(s3json.Key.Render(), 0, 0)

	zap.L().Debug("starting to work walk on", zap.String("url", s3json.URL().String()))
	go_for_a_walk(s3json)

	zap.L().Debug("walk done", zap.String("key", s3json.Key.Render()))

	return nil
}

package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	common "github.com/GSA-TTS/jemison/internal/common"
	kv "github.com/GSA-TTS/jemison/internal/kv"
	"github.com/GSA-TTS/jemison/internal/util"
	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pingcap/log"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

func host_and_path(job *river.Job[common.FetchArgs]) string {
	var u url.URL
	u.Scheme = job.Args.Scheme
	u.Host = job.Args.Host
	u.Path = job.Args.Path
	return u.String()
}

func chunkwiseSHA1(filename string) []byte {

	// Open the file for reading.
	tFile, err := os.Open(filename)
	if err != nil {
		zap.L().Error("could not open temp file for encoding to B64")
	}
	defer tFile.Close()
	// Compute the SHA1 going chunk-by-chunk
	h := sha1.New()
	reader := bufio.NewReader(tFile)
	// FIXME: make this a param in the config.
	chunkSize := 4 * 1024
	bytesRead := 0
	buf := make([]byte, chunkSize)
	for {
		n, err := reader.Read(buf)
		bytesRead += n

		if err != nil {
			if err != io.EOF {
				zap.L().Error("chunk error reading")
			}
			break
		}
		chunk := buf[0:n]
		// https://pkg.go.dev/crypto/sha1#example-New
		io.Writer.Write(h, chunk)
	}

	return h.Sum(nil)
}

func getUrlToFile(u url.URL) (string, int64, []byte) {
	getResponse, err := retryablehttp.Get(u.String())
	if err != nil {
		zap.L().Fatal("cannot GET content",
			zap.String("url", u.String()),
		)
	}
	zap.L().Debug("successful GET response")
	// Create a temporary file to download the HTML to.
	temporaryFilename := uuid.NewString()
	outFile, err := os.Create(temporaryFilename)
	if err != nil {
		zap.L().Panic("cannot create temporary file", zap.String("filename", temporaryFilename))
	}
	defer outFile.Close()

	// Copy the Get Reader to a file Writer
	// Should consume little/no RAM.
	// Destination, Source
	bytesRead, err := io.Copy(outFile, getResponse.Body)
	if err != nil {
		zap.L().Panic("could not copy GET to file",
			zap.String("url", u.String()),
			zap.String("filename", temporaryFilename))
	}
	getResponse.Body.Close()
	// Now, it is in a file.
	// Compute the SHA1
	theSHA := chunkwiseSHA1(temporaryFilename)
	return temporaryFilename, bytesRead, theSHA
}

// func getFilesize(filename string) int64 {
// 	fileInfo, err := os.Stat(filename)
// 	if err != nil {
// 		zap.L().Error("could not get filesize", zap.String("filename", filename))
// 	}
// 	return fileInfo.Size()
// }

func fetch_page_content(job *river.Job[common.FetchArgs]) (map[string]string, error) {
	u := url.URL{
		Scheme: job.Args.Scheme,
		Host:   job.Args.Host,
		Path:   job.Args.Path,
	}

	// This holds us up so that the parallel workers don't spam the host.
	common.BackoffLoop(job.Args.Host, polite_sleep, &last_hit, &last_backoff)

	zap.L().Debug("checking the hit cache")

	headResp, err := retryablehttp.Head(u.String())
	if err != nil {
		return nil, err
	}

	contentType := headResp.Header.Get("content-type")
	log.Debug("checking HEAD MIME type", zap.String("content-type", contentType))
	if !util.IsSearchableMimeType(contentType) {
		return nil, fmt.Errorf("non-indexable MIME type: %s", u.String())
	}

	// Write the raw content to a file.
	tempFilename, bytesRead, theSHA := getUrlToFile(u)
	key := util.CreateS3Key(util.ToScheme(job.Args.Scheme), job.Args.Host, job.Args.Path, util.Raw)

	defer func(u url.URL, key *util.Key) {
		err := os.Remove(tempFilename)
		if err != nil {
			zap.L().Error("could not remove temp file",
				zap.String("url", u.String()),
				zap.String("key", key.Render()))
		}
	}(u, key)

	// Stream that file over to S3
	s3 := kv.NewS3(ThisServiceName)
	s3.FileToS3(key, tempFilename, util.GetMimeType(contentType))

	response := map[string]string{
		"raw":            key.Render(),
		"sha1":           fmt.Sprintf("%x", theSHA),
		"content-length": fmt.Sprintf("%d", bytesRead),
		"scheme":         job.Args.Scheme,
		"host":           job.Args.Host,
		"path":           job.Args.Path,
	}
	// FIXME
	// There is a texinfo standard library for normalizing content types.
	// Consider using it. I want a simplified string, not utf-8 etc.
	response["content-type"] = util.GetMimeType(response["content-type"])

	zap.L().Debug("content read",
		zap.String("content-length", response["content-length"]),
	)

	// Copy in all of the response headers.
	// This used to be the GET headers, but... they're hiding.
	// Going to do this for now, because I don't know what we'll need.
	for k := range headResp.Header {
		response[strings.ToLower(k)] = headResp.Header.Get(k)
	}

	return response, nil
}

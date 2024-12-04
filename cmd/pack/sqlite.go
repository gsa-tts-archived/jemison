package main

import (
	"log"
	"strings"
	"syscall"

	"github.com/GSA-TTS/jemison/internal/kv"
	"github.com/GSA-TTS/jemison/internal/sqlite"
	"go.uber.org/zap"
)

type Ping struct{}

var limit = make(chan Ping, 1)

func getSpaceAvailable() uint64 {
	var stat syscall.Statfs_t

	err := syscall.Statfs("/", &stat)
	if err != nil {
		zap.L().Fatal("could not get disk space for packing")
	}

	// Available space in bytes
	availableBytes := stat.Bavail * uint64(stat.Bsize)
	return availableBytes
}

func optionalSlash(path string) string {
	if strings.HasSuffix(path, "/") {
		return path
	} else {
		return path + "/"
	}
}

func PackSqliteDatabase(host string) {

	// Look into the "extract" bucket and get a
	// list of all the objects there.
	s3 := kv.NewS3("extract")

	list_of_objects, err := s3.List(optionalSlash(host))
	if err != nil {
		zap.L().Fatal("could not list objects",
			zap.String("bucket", "extract"),
			zap.String("host", host))
	}

	// Calculate size
	var size int64 = 0
	for _, s3obj := range list_of_objects {
		//zap.L().Debug("size", zap.String("key", s3obj.Key), zap.Int64("size", s3obj.Size))
		size += s3obj.Size
	}
	available := getSpaceAvailable()
	if (float64(size) * 1.1) > float64(available) {
		// If we have more to pack than is available on the machine
		// 1) we're in trouble
		// 2) we want to send this back to the queue.
		//
		// FIXME: Add a client to queue this back for ourselves.
		zap.L().Warn("not enough space to pack the database",
			zap.Float64("database size", float64(size)), zap.Float64("available", float64(available)))
		return
	}

	pt, err := sqlite.CreatePackTable(sqlite.SqliteFilename(host))
	if err != nil {
		log.Println("Could not create pack table for", host)
		log.Fatal(err)
	}

	for _, s3obj := range list_of_objects {
		s3json, err := s3.S3PathToS3JSON(s3obj.Key)
		if err != nil {
			zap.L().Error("could not fetch object for packing",
				zap.String("key", s3json.Key.Render()))
		}

		contentType := s3json.GetString("content-type")
		switch contentType {
		case "text/html":
			packHtml(pt, s3json)
		case "application/pdf":
			packPdf(pt, s3json)
		}
	}

	// Cleanup time.
	pt.DB.Close()
	pt.PrepForNetwork()
}

func PackPostgresDatabase(host string) {

	// Look into the "extract" bucket and get a
	// list of all the objects there.
	s3 := kv.NewS3("extract")

	list_of_objects, err := s3.List(optionalSlash(host))
	if err != nil {
		zap.L().Fatal("could not list objects",
			zap.String("bucket", "extract"),
			zap.String("host", host))
	}

	for _, s3obj := range list_of_objects {
		s3json, err := s3.S3PathToS3JSON(s3obj.Key)
		if err != nil {
			zap.L().Error("could not fetch object for packing",
				zap.String("key", s3json.Key.Render()))
		}

		contentType := s3json.GetString("content-type")
		switch contentType {
		case "text/html":
			packPgHtml(SDB0, s3json)
		case "application/pdf":
			packPgPdf(SDB0, s3json)
		}
	}

}

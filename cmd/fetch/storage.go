package main

import "github.com/GSA-TTS/jemison/internal/kv"

var fetchStorage *kv.S3JSON

func InitializeStorage() {
	fetchStorage = kv.NewS3JSON("fetch")
}

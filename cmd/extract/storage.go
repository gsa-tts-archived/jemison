package main

import "github.com/GSA-TTS/jemison/internal/kv"

var fetchStorage kv.S3
var extractStorage kv.S3

func InitializeStorage() {
	fetchStorage = kv.NewKV("fetch")
	extractStorage = kv.NewKV("extract")
}

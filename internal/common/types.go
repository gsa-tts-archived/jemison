package common

import "net/http"

type ExtractArgs struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Path   string `json:"path"`
}

func (ExtractArgs) Kind() string {
	return "extract"
}

type FetchArgs struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Path   string `json:"path"`
}

func (FetchArgs) Kind() string {
	return "fetch"
}

type PackArgs struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Path   string `json:"path"`
}

func (PackArgs) Kind() string {
	return "pack"
}

type ServeArgs struct {
	Filename string `json:"filename"`
}

func (ServeArgs) Kind() string {
	return "serve"
}

type WalkArgs struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Path   string `json:"path"`
}

func (WalkArgs) Kind() string {
	return "walk"
}

type HttpResponse func(w http.ResponseWriter, r *http.Request)

// VALIDATOR TYPES
var ValidateFetchQueue = "validate_fetch"

type ValidateFetchArgs struct {
	Fetch FetchArgs
}

func (ValidateFetchArgs) Kind() string {
	return ValidateFetchQueue
}

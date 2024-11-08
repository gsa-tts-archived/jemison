package sitemaps

import "net/url"

type BaseSitemap struct {
	URLs []url.URL
}
type RemoteSitemap struct {
	Scheme string
	Host   string
	Path   string
	BaseSitemap
}

type FileSitemap struct {
	Filename string
	BaseSitemap
}

func NewFileSitemap(filename string) FileSitemap {
	return FileSitemap{Filename: filename}
}

func (fsm *FileSitemap) Load() {

}

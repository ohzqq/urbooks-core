package urbooks

import (
	"fmt"
	"net/url"

	"github.com/ohzqq/urbooks-core/calibredb"
)

var _ = fmt.Sprintf("%v", "")

func Libraries() []string {
	return cfg.list
}

func Lib(l string) *Library {
	return cfg.libs[l]
}

func DefaultLib() *Library {
	return Lib(Cfg().Opts["default"])
}

type Library struct {
	Cfg            *libCfg
	Name           string
	Path           string
	DefaultRequest *Request
	DB             *calibredb.Lib
	req            *Request
}

func NewLibrary(name, path string) *Library {
	l := Library{}
	l.Cfg = Cfg().libCfg[name]
	l.Name = name
	l.Path = path
	l.DB = calibredb.NewLib(path)
	return &l
}

func (l *Library) IsAudiobooks() bool {
	return l.Cfg.Audiobooks
}

func (l *Library) NewBook() *Book {
	book := NewBook()
	var q = make(url.Values)
	q.Set("library", l.Name)
	book.query = q
	return book
}

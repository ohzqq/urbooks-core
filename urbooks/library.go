package urbooks

import (
	"github.com/ohzqq/urbooks-core/calibredb"
	"golang.org/x/exp/slices"
)

func Libraries() []string {
	return cfg.list
}

func Lib(l string) *Library {
	return cfg.libs[l]
}

func DefaultLib() *Library {
	var lib = Libraries()[0]
	if l := Cfg().Opts["default"]; l != "" {
		if slices.Contains(Libraries(), l) {
			lib = l
		}
	}
	return Lib(lib)
}

type Library struct {
	Cfg            *libCfg
	Name           string
	Path           string
	DefaultRequest *Request
	DB             *calibredb.Lib
	pref           *calibredb.Preferences
	req            *Request
}

func NewLibrary(name, path string) *Library {
	l := Library{}
	l.Cfg = Cfg().libCfg[name]
	l.Name = name
	l.Path = path
	l.DB = calibredb.NewLib(path)
	l.pref = l.DB.Preferences
	return &l
}

func (l *Library) IsAudiobooks() bool {
	return l.Cfg.Audiobooks
}

//func (l *Library) NewBook() *Book {
//  book := NewBook()
//  var q = make(url.Values)
//  q.Set("library", l.Name)
//  book.query = q
//  return book
//}

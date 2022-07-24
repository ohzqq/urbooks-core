package urbooks

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

func Cfg() *config {
	return cfg
}

type libCfg struct {
	Name       string
	Path       string
	Default    bool
	Audiobooks bool
	WebOpts    webOpts `mapstructure:"website_options"`
}

func newLibCfg(v *viper.Viper, name string) *libCfg {
	cfg := &libCfg{Name: name}
	err := v.Sub(name).Unmarshal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

type webOpts struct {
	Path  string
	URL   string
	Files string
	Feeds map[string][]string
}

type config struct {
	libs   map[string]*Library
	Opts   map[string]string
	list   []string
	libCfg map[string]*libCfg
}

func (c *config) CatMin() int {
	min, err := strconv.Atoi(c.Opts["cat_min"])
	if err != nil {
		return 1
	}
	return min
}

var cfg = &config{
	libCfg: make(map[string]*libCfg),
	libs:   make(map[string]*Library),
}

func InitConfig(opts map[string]string) {
	cfg.Opts = opts
}

func InitLibraries(v *viper.Viper, web bool) {
	libs := v.AllSettings()
	if len(libs) == 0 {
		log.Fatal("no libraries list in config")
	}
	for lib, _ := range libs {
		Cfg().list = append(Cfg().list, lib)

		libcfg := newLibCfg(v, lib)

		var libPath string
		switch web {
		case true:
			libPath = filepath.Join(libcfg.WebOpts.Path, lib)
		case false:
			libPath = filepath.Join(libcfg.Path, lib)
		}

		Cfg().libCfg[lib] = libcfg

		if _, err := os.Stat(libPath); errors.Is(err, os.ErrNotExist) {
			log.Fatalf("%v does not exist or cannot be found at %v, check the path in the config", libcfg.Name, libPath)
		}

		newLib := NewLibrary(lib, libPath)
		newLib.ConnectDB()
		newLib.DefaultRequest = NewRequest(lib).From("books")

		if sort := Cfg().Opts["sort"]; sort != "" {
			newLib.DefaultRequest.Sort(sort)
		}

		if order := Cfg().Opts["order"]; order != "" {
			newLib.DefaultRequest.Set("order", order)
		}

		Cfg().libs[lib] = newLib

		newLib.GetDBPreferences()
	}
}

func parseHiddenCategories(d json.RawMessage) []string {
	var hidden []string
	err := json.Unmarshal(d, &hidden)
	if err != nil {
		log.Fatal(err)
	}
	return hidden
}

func parseSavedSearches(d json.RawMessage) map[string]string {
	var searches map[string]string
	err := json.Unmarshal(d, &searches)
	if err != nil {
		log.Fatal(err)
	}
	return searches
}

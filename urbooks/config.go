package urbooks

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

var _ = fmt.Sprintf("%v", "")

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

func InitLibraries(v *viper.Viper, libs map[string]string, web bool) {
	if len(libs) == 0 {
		log.Fatal("no libraries list in config")
	}
	for lib, _ := range libs {
		Cfg().list = append(Cfg().list, lib)

		libcfg := libCfg{}
		libcfg.Name = lib
		err := v.Sub(lib).Unmarshal(&libcfg)
		if err != nil {
			log.Fatal(err)
		}
		var libPath string
		switch web {
		case true:
			libPath = filepath.Join(libcfg.WebOpts.Path, lib)
		case false:
			libPath = filepath.Join(libcfg.Path, lib)
		}

		Cfg().libCfg[lib] = &libcfg

		if _, err := os.Stat(libPath); errors.Is(err, os.ErrNotExist) {
			log.Fatalf("%v does not exist or cannot be found at %v, check the path in the config", libcfg.Name, libPath)
		}

		newLib := NewLibrary(lib, libPath)
		newLib.DefaultRequest = NewRequest(lib).From("books")

		if sort := Cfg().Opts["sort"]; sort != "" {
			newLib.DefaultRequest.Sort(sort)
		}

		if order := Cfg().Opts["order"]; order != "" {
			newLib.DefaultRequest.Set("order", order)
		}

		Cfg().libs[lib] = newLib
	}
}

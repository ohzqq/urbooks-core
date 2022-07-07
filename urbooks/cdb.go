package urbooks

import (
	"log"
	"net"
	"net/url"
	"regexp"
	"time"

	"github.com/ohzqq/avtools/avtools"
	"golang.org/x/exp/slices"
)

type CalibreUserCfg struct {
	Url      string
	Username string
	Password string
}

type cdbCmd struct {
	lib   string
	user  CalibreUserCfg
	input string
	media *avtools.Media
	cmd   string
	args  []string
}

func NewCalibredbCmd() *cdbCmd {
	return &cdbCmd{}
}

func (c *cdbCmd) SetUser(config CalibreUserCfg) *cdbCmd {
	c.user = config
	return c
}

func (c *cdbCmd) SetLib(l string) *cdbCmd {
	if !slices.Contains(Libraries(), l) {
		log.Fatal("this is not a library")
	}
	c.lib = l
	return c
}

func (c *cdbCmd) Add(input string) *cdbCmd {
	c.cmd = "add"
	c.input = input
	c.media = avtools.NewMedia(input).JsonMeta().Unmarshal()
	return c
}

func (c *cdbCmd) MediaMetaToBook() *Book {
	//fmt.Printf("media meta: %T\n", NewBookMeta(c.media.Meta.Tags))
	book := NewBook(c.lib)
	titleRegex := regexp.MustCompile(`(?P<title>.*) \[(?P<series>.*), Book (?P<position>.*)\]$`)
	titleAndSeries := titleRegex.FindStringSubmatch(c.media.Meta.Tags["title"])
	book.NewColumn("title").SetValue(titleAndSeries[titleRegex.SubexpIndex("title")])
	book.NewItem("series").
		SetValue(titleAndSeries[titleRegex.SubexpIndex("series")]).
		Set("position", titleAndSeries[titleRegex.SubexpIndex("position")])
	book.NewCategory("authors").Split(c.media.Meta.Tags["artist"], true)
	book.NewCategory("narrators").Split(c.media.Meta.Tags["composer"], true)
	book.NewColumn("description").SetValue(c.media.Meta.Tags["comment"])
	book.NewCategory("tags").Split(c.media.Meta.Tags["genre"], false)
	return book
}

func (c *cdbCmd) ParseCfg() *cdbCmd {
	if l := c.user.Url; l != "" {
		if c.user.IsOnline() {
			c.args = append(c.args, "with-library")
			u, err := url.Parse(l)
			if err != nil {
				log.Fatal(err)
			}

			u.Fragment = c.lib
			c.args = append(c.args, u.String())
		}
	}

	if p := c.user.Password; p != "" {
		c.args = append(c.args, "password")
		c.args = append(c.args, p)
	}

	if u := c.user.Username; u != "" {
		c.args = append(c.args, "username")
		c.args = append(c.args, u)
	}

	if !slices.Contains(c.args, "with-library") {
		c.args = append(c.args, "with-library")
		c.args = append(c.args, Lib(c.lib).Path)
	}
	return c
}

func (c CalibreUserCfg) IsOnline() bool {
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", c.Url, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

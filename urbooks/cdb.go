package urbooks

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/ohzqq/avtools/avtools"
	"github.com/ohzqq/urbooks-core/calibredb"
	"golang.org/x/exp/slices"
)

type CalibreServerCfg struct {
	Url      string
	Username string
	Password string
	url      *url.URL
}

type cdbCmd struct {
	lib      *Library
	server   CalibreServerCfg
	input    string
	verbose  bool
	media    *avtools.Media
	book     *Book
	cdbCmd   string
	localCmd string
	tmp      *os.File
	args     []string
	cmd      *exec.Cmd
}

func NewCalibredbCmd() *cdbCmd {
	return &cdbCmd{}
}

func (c *cdbCmd) Verbose(v bool) *cdbCmd {
	c.verbose = v
	return c
}

func (c *cdbCmd) Run() string {
	//fmt.Printf("cdbCmd: %v\n", c.cdbCmd)
	//fmt.Printf("localCmd: %v\n", c.localCmd)
	//fmt.Printf("c.lib: %v\n", c.lib)
	switch {
	case c.cdbCmd != "":
		return c.runExec()
	default:
		return c.runLocal()
	}
	return ""
}

func (c *cdbCmd) runLocal() string {
	if c.verbose {
		fmt.Printf("local command: %v\n", c.localCmd)
	}
	switch c.localCmd {
	case "list libraries":
		return strings.Join(Libraries(), ", ")
	case "list fields":
		return strings.Join(c.lib.DB.AllFields(), ", ")
	}
	return ""
}

func (c *cdbCmd) runExec() string {
	c.buildCmd()

	if c.tmp != nil {
		defer os.Remove(c.tmp.Name())
	}

	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	c.cmd.Stderr = &stderr
	c.cmd.Stdout = &stdout

	err := c.cmd.Run()

	fmt.Printf("%v\n", c.cmd.String())
	if err != nil {
		fmt.Printf("%v\n", c.cmd.String())
		log.Fatalf("%v finished with error: %v\n", c.cdbCmd, stderr.String())
	}

	if len(stdout.Bytes()) > 0 {
		return stdout.String()
	}

	if c.verbose {
		fmt.Println(c.cmd.String())
	}

	return ""
}

func (c *cdbCmd) buildCmd() *cdbCmd {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, c.cdbCmd)
	cmdArgs = append(cmdArgs, c.ParseCfg()...)
	cmdArgs = append(cmdArgs, c.args...)
	c.cmd = exec.Command("calibredb", cmdArgs...)
	return c
}

func (c *cdbCmd) appendArgs(arg ...string) *cdbCmd {
	c.args = append(c.args, arg...)
	return c
}

func (c *cdbCmd) SetServer(config CalibreServerCfg) *cdbCmd {
	c.server = config
	return c
}

func (c *cdbCmd) SetLib(l string) *cdbCmd {
	switch {
	case l == "":
		c.lib = DefaultLib()
	case !slices.Contains(Libraries(), l):
		log.Fatal("this is not a library")
	default:
		c.lib = Lib(l)
	}
	return c
}

func (c *cdbCmd) List(arg string) *cdbCmd {
	switch arg {
	case "fields":
		c.localCmd = "list fields"
		//c.setMetadataCmd()
		c.appendArgs("--list-fields")
	case "libs":
		c.localCmd = "list libraries"
	}
	return c
}

func (c *cdbCmd) Add(input, cover string) *cdbCmd {
	id, err := c.addBook(input, cover)
	if err != nil {
		log.Fatal(err)
	}
	metaCmd := &cdbCmd{
		media:   c.media,
		lib:     c.lib,
		verbose: c.verbose,
		server:  c.server,
		book:    c.MediaMetaToBook(),
	}
	metaCmd.setMetadataCmd(id).Run()
	return c
}

func (c *cdbCmd) addBook(input, cover string) (string, error) {
	c.media = avtools.NewMedia(input).JsonMeta().Unmarshal()
	c.cdbCmd = "add"
	c.appendArgs("-d")
	if cover != "" {
		c.appendArgs("-c", cover)
	}
	c.appendArgs(input)
	if id := strings.Split(c.Run(), ": "); len(id) == 2 {
		return id[1], nil
	}
	return "", fmt.Errorf("import unsucessful")
}

func (c *cdbCmd) setMetadataCmd(id string) *cdbCmd {
	c.cdbCmd = "set_metadata"
	c.appendArgs(id)
	for field, value := range c.book.StringMap() {
		f := calibredb.GetCalibreField(field) + ":"
		switch {
		case field == "library":
		default:
			c.appendArgs("-f", f+value)
		}
	}
	return c
}

func shEscape(s string) string {
	escSpace := regexp.MustCompile(`(\s)`)
	return escSpace.ReplaceAllString(s, "\\$1")
}

func (c *cdbCmd) listCmd() *cdbCmd {
	c.cdbCmd = "list"
	return c
}

func (c *cdbCmd) MediaMetaToBook() *Book {
	//fmt.Printf("media meta: %T\n", NewBookMeta(c.media.Meta.Tags))
	book := NewBook(c.lib.Name)
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

func (c *cdbCmd) ParseCfg() []string {
	var args []string
	if l := c.server.Url; l != "" {
		u, err := url.Parse(l)
		if err != nil {
			log.Fatal(err)
		}
		u.Fragment = c.lib.Name
		c.server.url = u

		if c.server.IsOnline() {
			args = append(args, "--with-library")
			args = append(args, "'"+c.server.url.String()+"'")
		}
	}

	if p := c.server.Password; p != "" {
		args = append(args, "--password")
		args = append(args, "'"+p+"'")
	}

	if u := c.server.Username; u != "" {
		args = append(args, "--username")
		args = append(args, u)
	}

	if !slices.Contains(args, "--with-library") {
		args = append(args, "--with-library")
		switch c.lib {
		//case "-":
		//args = append(args, c.lib)
		default:
			args = append(args, c.lib.Path)
		}
	}

	return args
}

func (c CalibreServerCfg) IsOnline() bool {
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", c.url.Host, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

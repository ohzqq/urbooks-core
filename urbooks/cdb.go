package urbooks

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ohzqq/avtools/avtools"
	"github.com/ohzqq/urbooks-core/book"
	"github.com/ohzqq/urbooks-core/calibredb"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

type calibreCfg struct {
	cli *viper.Viper
	srv *viper.Viper
	url *url.URL
}

var cdb = calibreCfg{}

func CfgCdb(v *viper.Viper) {
	cdb.cli = v.Sub("cdb")
	cdb.srv = v
}

type cdbCmd struct {
	lib      *Library
	input    string
	verbose  bool
	media    *avtools.Media
	book     *book.Book
	cdbCmd   string
	localCmd string
	tmp      *os.File
	args     []string
	cmd      *exec.Cmd
}

type importCmd struct {
	*cdbCmd
	meta map[string]string
}

func NewCalibredbCmd() *cdbCmd {
	return &cdbCmd{}
}

func (c *cdbCmd) Verbose(v bool) *cdbCmd {
	c.verbose = v
	return c
}

func (c *cdbCmd) Run() (string, error) {
	switch {
	case c.cdbCmd != "":
		return c.runCdb()
	default:
		if c.verbose {
			fmt.Printf("local command: %v\n", c.localCmd)
		}
		switch c.localCmd {
		case "list libraries":
			return strings.Join(Libraries(), ", "), nil
		case "list fields":
			return strings.Join(c.lib.DB.AllFields(), ", "), nil
		}
	}
	return "", nil
}

func (c *cdbCmd) runCdb() (string, error) {
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

	if err != nil {
		fmt.Printf("%v\n", c.cmd.String())
		log.Fatalf("%v finished with error: %v\n", c.cdbCmd, stderr.String())
		return "", fmt.Errorf("%v\n", stderr.String())
	}

	if len(stdout.Bytes()) > 0 {
		return stdout.String(), nil
	}

	if c.verbose {
		fmt.Println(c.cmd.String())
	}

	return "", nil
}

type CalibreSearch struct {
	Authors   string
	Added     string
	Limit     string
	Narrators string
	Order     string
	Published string
	Publisher string
	Rating    string
	Series    string
	Sort      string
	Tags      string
	Title     string
}

func (c *cdbCmd) searchCmd() *cdbCmd {
	c.setCdbCmd("search")
	return c
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

func (c *cdbCmd) WithLib(l string) *cdbCmd {
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
		c.appendArgs("--list-fields")
	case "libs":
		c.localCmd = "list libraries"
	}
	return c
}

func (c *cdbCmd) Import(input, cover, metaFile string) *cdbCmd {
	id, err := c.addBook(input, cover)
	if err != nil {
		log.Fatal(err)
	}

	meta := make(map[string]string)
	if metaFile != "" {
		file, err := os.Open(metaFile)
		if err != nil {
			log.Fatal(err)
		}
		_, err = toml.NewDecoder(file).Decode(&meta)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		b := book.MediaMetaToBook(c.lib.Name, c.media)
		meta = b.StringMap(false)
	}
	delete(meta, "id")

	metaCmd := &cdbCmd{
		media:   c.media,
		lib:     c.lib,
		verbose: c.verbose,
	}
	_, err = metaCmd.setMetadataCmd(id, meta).Run()
	switch {
	case err != nil:
		log.Fatal(err)
	case err == nil:
		fmt.Printf("imported %v, id %v\n", c.media.GetTag("title"), id)
	}
	return c
}

func (c *cdbCmd) Export(ids, dir, fmts string) *cdbCmd {
	c.setCdbCmd("export")
	c.appendArgs(ids)

	if c.lib.IsAudiobooks() {
		c.appendArgs("--dont-update-metadata")
	}

	if dir != "" {
		c.appendArgs("--to-dir", dir)
	}

	if fmts != "" {
		c.appendArgs("--formats", fmts)
	}

	_, err := c.Run()
	switch {
	case err != nil:
		log.Fatal(err)
	case err == nil:
		fmt.Println("export successful")
	}
	return c
}

func (c *cdbCmd) Remove(id string) *cdbCmd {
	c.setCdbCmd("remove").appendArgs(id).Run()
	return c
}

func (c *cdbCmd) setCdbCmd(cmd string) *cdbCmd {
	c.cdbCmd = cmd
	if cdb.cli.IsSet(c.cdbCmd) {
		for _, o := range cdb.cli.GetStringSlice(c.cdbCmd) {
			c.appendArgs(o)
		}
	}
	return c
}

func (c *cdbCmd) addBook(input, cover string) (string, error) {
	c.media = avtools.NewMedia(input).JsonMeta().Unmarshal()

	c.setCdbCmd("add")

	if cover != "" {
		c.appendArgs("-c", cover)
	}

	c.appendArgs(input)

	str, err := c.Run()
	if err != nil {
		log.Fatal(err)
	}

	if id := strings.Split(str, ": "); len(id) == 2 {
		return id[1], nil
	}

	return "", fmt.Errorf("import unsucessful")
}

func (c *cdbCmd) setMetadataCmd(id string, meta map[string]string) *cdbCmd {
	c.setCdbCmd("set_metadata")
	c.appendArgs(id)

	for field, val := range meta {
		f := calibredb.GetCalibreField(field) + ":"
		if strings.HasPrefix(field, "#") || slices.Contains(book.EditableFields, field) {
			c.appendArgs("-f", f+val)
		}
	}

	return c
}

func (c *cdbCmd) listCmd() *cdbCmd {
	return c.setCdbCmd("list")
}

func (c *cdbCmd) ParseCfg() []string {
	var args []string
	if cdb.srv.IsSet("url") {
		u, err := url.Parse(cdb.srv.GetString("url"))
		if err != nil {
			log.Fatal(err)
		}
		u.Fragment = c.lib.Name
		cdb.url = u

		if cdb.IsOnline() {
			args = append(args, "--with-library")
			args = append(args, "'"+cdb.url.String()+"'")
		}
	}

	if cdb.srv.IsSet("password") {
		args = append(args, "--password")
		args = append(args, "'"+cdb.srv.GetString("password")+"'")
	}

	if cdb.srv.IsSet("username") {
		args = append(args, "--username")
		args = append(args, cdb.srv.GetString("username"))
	}

	if !slices.Contains(args, "--with-library") {
		args = append(args, "--with-library")
		args = append(args, c.lib.Path)
	}

	return args
}

func (c calibreCfg) IsOnline() bool {
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", cdb.url.Host, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

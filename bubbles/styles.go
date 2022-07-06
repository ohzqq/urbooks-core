package urbooksTui

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

type TuiTheme struct {
	Default_fg string
	Default_bg string
	Black      string
	Blue       string
	Cyan       string
	Green      string
	Grey       string
	Pink       string
	Purple     string
	Red        string
	White      string
	Yellow     string
}

var Theme = TuiTheme{
	Black:  "#000000",
	Blue:   "#0000FF",
	Cyan:   "#00FFFF",
	Green:  "#00FF00",
	Grey:   "#C0C0C0",
	Pink:   "#FF00FF",
	Purple: "#800080",
	Red:    "#FF0000",
	White:  "#FFFFFF",
	Yellow: "#FFFF00",
}

func Config(v *viper.Viper) {
	err := v.UnmarshalKey("theme", &Theme)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}
}

const (
	bullet   = "•"
	ellipsis = "…"
)

type ItemStyle struct {
	NormalItem   lipgloss.Style
	CurrentItem  lipgloss.Style
	SelectedItem lipgloss.Style
}

func ItemStyles() (s ItemStyle) {
	s.NormalItem = lipgloss.NewStyle().Foreground(lipgloss.Color(Theme.Default_fg)).Margin(0, 1, 0, 2)
	s.CurrentItem = lipgloss.NewStyle().Foreground(lipgloss.Color(Theme.Green)).Reverse(true).Margin(0, 1, 0, 2)
	s.SelectedItem = lipgloss.NewStyle().Foreground(lipgloss.Color(Theme.Grey)).Margin(0, 1, 0, 2)
	return s
}

func ListStyles() (s list.Styles) {
	verySubduedColor := lipgloss.AdaptiveColor{Light: Theme.White, Dark: Theme.Grey}
	subduedColor := lipgloss.AdaptiveColor{Light: Theme.Grey, Dark: Theme.White}

	s.TitleBar = lipgloss.NewStyle().Padding(0, 0, 0, 2)

	s.Title = lipgloss.NewStyle().
		Background(lipgloss.Color(Theme.Purple)).
		Foreground(lipgloss.Color(Theme.Black)).
		Padding(0, 1)

	s.Spinner = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: Theme.Black, Dark: Theme.Cyan})

	s.FilterPrompt = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: Theme.Black, Dark: Theme.Pink})

	s.FilterCursor = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: Theme.Black, Dark: Theme.Yellow})

	s.DefaultFilterCharacterMatch = lipgloss.NewStyle().Underline(true)

	s.StatusBar = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: Theme.Black, Dark: Theme.Blue}).
		Padding(0, 0, 1, 2)

	s.StatusEmpty = lipgloss.NewStyle().Foreground(subduedColor)

	s.StatusBarActiveFilter = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: Theme.Black, Dark: Theme.Purple})

	s.StatusBarFilterCount = lipgloss.NewStyle().Foreground(verySubduedColor)

	s.NoItems = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: Theme.Black, Dark: Theme.Grey})

	s.ArabicPagination = lipgloss.NewStyle().Foreground(subduedColor)

	s.PaginationStyle = lipgloss.NewStyle().PaddingLeft(2) //nolint:gomnd

	s.HelpStyle = lipgloss.NewStyle().Padding(1, 0, 0, 2)

	s.ActivePaginationDot = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: Theme.Black, Dark: Theme.Pink}).
		SetString(bullet)

	s.InactivePaginationDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		SetString(bullet)

	s.DividerDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		SetString(" " + bullet + " ")

	return s
}

func helpStyles() help.Styles {
	styles := help.Styles{}

	keyStyle := lipgloss.NewStyle().PaddingRight(1).Foreground(lipgloss.Color(Theme.Green))
	descStyle := lipgloss.NewStyle().PaddingRight(1).Foreground(lipgloss.Color(Theme.Blue))
	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(Theme.Pink))

	styles.ShortKey = keyStyle
	styles.ShortDesc = descStyle
	styles.ShortSeparator = sepStyle
	styles.FullKey = keyStyle.Copy()
	styles.FullDesc = descStyle.Copy()
	styles.FullSeparator = sepStyle.Copy()
	styles.Ellipsis = sepStyle.Copy()

	return styles
}

func RenderMarkdown(md string) string {
	var (
		metaStyle = template.Must(template.New("mdStyle").Parse(styleTmpl))
		style     bytes.Buffer
	)

	err := metaStyle.Execute(&style, Theme)
	if err != nil {
		log.Fatal(err)
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes(style.Bytes()),
		glamour.WithPreservedNewLines(),
	)
	if err != nil {
		log.Fatal(err)
	}

	str, err := renderer.Render(md)
	if err != nil {
		log.Fatal(err)
	}

	return str
}

var styleTmpl = `
{
  "document": {
    "block_prefix": "",
    "block_suffix": "",
    "color": "{{.Default_fg}}",
    "background_color": "{{.Default_bg}}",
    "margin": 0
  },
  "block_quote": {
    "indent": 1,
    "indent_token": "│ "
  },
  "paragraph": {
    "block_suffix": ""
  },
  "list": {
    "level_indent": 2
  },
  "heading": {
    "block_suffix": "",
    "color": "{{.Pink}}",
    "bold": true
  },
  "h1": {
    "prefix": " ",
    "suffix": " ",
    "color": "{{.Default_bg}}",
    "background_color": "{{.Blue}}",
    "bold": true
  },
  "h2": {
    "prefix": "## "
  },
  "h3": {
    "prefix": "### "
  },
  "h4": {
    "prefix": "#### "
  },
  "h5": {
    "prefix": "##### "
  },
  "h6": {
    "prefix": "###### ",
    "bold": false
  },
  "text": {},
  "strikethrough": {
    "crossed_out": true
  },
  "emph": {
    "italic": true
  },
  "strong": {
    "color": "{{.Cyan}}",
    "bold": true
  },
  "hr": {
    "color": "{{.Pink}}",
    "format": "\n--------\n"
  },
  "item": {
    "block_prefix": "• "
  },
  "enumeration": {
    "block_prefix": ". "
  },
  "html_block": {},
  "html_span": {}
}
`

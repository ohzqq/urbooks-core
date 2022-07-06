package bubbles

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/statusbar"
)

type headerBar struct {
	model   statusbar.Bubble
	width   int
	height  int
	focused bool
}

func newHeader() *headerBar {
	col1Style := statusbar.ColorConfig{
		Foreground: lipgloss.AdaptiveColor{
			Light: Theme.Cyan,
			Dark:  Theme.Cyan,
		},
		Background: lipgloss.AdaptiveColor{
			Light: Theme.Default_bg,
			Dark:  Theme.Default_bg,
		},
	}
	col2Style := statusbar.ColorConfig{
		Foreground: lipgloss.AdaptiveColor{
			Light: Theme.Green,
			Dark:  Theme.Green,
		},
		Background: lipgloss.AdaptiveColor{
			Light: Theme.Default_bg,
			Dark:  Theme.Default_bg,
		},
	}
	col3Style := statusbar.ColorConfig{
		Foreground: lipgloss.AdaptiveColor{
			Light: Theme.Purple,
			Dark:  Theme.Purple,
		},
		Background: lipgloss.AdaptiveColor{
			Light: Theme.Default_bg,
			Dark:  Theme.Default_bg,
		},
	}
	col4Style := statusbar.ColorConfig{
		Foreground: lipgloss.AdaptiveColor{
			Light: Theme.Pink,
			Dark:  Theme.Pink,
		},
		Background: lipgloss.AdaptiveColor{
			Light: Theme.Default_bg,
			Dark:  Theme.Default_bg,
		},
	}
	width := TermWidth()
	status := statusbar.New(col1Style, col2Style, col3Style, col4Style)
	status.FirstColumn = "URbooks"
	status.Width = width
	return &headerBar{
		model:   status,
		focused: false,
		width:   width,
		height:  1,
	}
}

type statusBarUpdateMsg string

func (h *headerBar) setCol1(text string) tea.Cmd {
	h.model.FirstColumn = text
	return func() tea.Msg {
		return statusBarUpdateMsg(text)
	}
}

func (h *headerBar) setCol2(text string) tea.Cmd {
	h.model.SecondColumn = text
	return func() tea.Msg {
		return statusBarUpdateMsg(text)
	}
}

func (h *headerBar) setCol3(text string) tea.Cmd {
	h.model.ThirdColumn = text
	return func() tea.Msg {
		return statusBarUpdateMsg(text)
	}
}

func (h *headerBar) setCol4(text string) tea.Cmd {
	h.model.FourthColumn = text
	return func() tea.Msg {
		return statusBarUpdateMsg(text)
	}
}

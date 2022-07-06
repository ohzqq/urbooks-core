package bubbles

import (
	//"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyContextMenu struct {
	title   string
	model   help.Model
	keyMap  KeyMap
	width   int
	height  int
	focused bool
}

func NewKeyContextMenu() *KeyContextMenu {
	menu := KeyContextMenu{}
	menu.focused = false
	return &menu
}

func (k *KeyContextMenu) SetKeyMap(keys KeyMap) *KeyContextMenu {
	k.keyMap = keys
	return k
}

func (k *KeyContextMenu) SetTitle(t string) *KeyContextMenu {
	k.title = t
	return k
}

func (k *KeyContextMenu) SetWidth(w int) *KeyContextMenu {
	k.width = w
	return k
}

func (k *KeyContextMenu) SetFullHelp(keys [][]key.Binding) *KeyContextMenu {
	//k.FullHelpView(keys)
	return k
}

func (k *KeyContextMenu) Model() *KeyContextMenu {
	menu := help.New()
	menu.ShowAll = true
	menu.Styles = helpStyles()
	if k.width != 0 {
		menu.Width = k.width
	}
	return k
}

func (k *KeyContextMenu) SetFocused(f bool) {
	k.focused = f
}

type newSortMsg struct{}

func (k *KeyContextMenu) sorted() tea.Cmd {
	return func() tea.Msg {
		return newSortMsg{}
	}
}

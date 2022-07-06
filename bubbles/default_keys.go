package bubbles

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap map[string]key.Binding

func (s KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		s["ExitScreen"],
		s["FullHelp"],
		s["Enter"],
	}
}

func (s KeyMap) FullHelp() [][]key.Binding {
	var keys [][]key.Binding
	first := []key.Binding{
		s["Quit"],
		s["ExitScreen"],
		s["FullHelp"],
		s["PrevMenu"],
	}
	keys = append(keys, first)
	second := []key.Binding{
		s["SortBy"],
		s["ChangeLibrary"],
	}
	keys = append(keys, second)
	third := []key.Binding{
		s["ToggleItem"],
		s["DeselectAll"],
		s["SelectAll"],
	}
	keys = append(keys, third)
	return keys
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		"Quit": key.NewBinding(
			key.WithKeys("ctrl+c", "esc"),
			key.WithHelp("ctrl+c", "quit"),
		),
		"Enter": key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select item"),
		),
		"FullHelp": key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "full help"),
		),
		"MetaViewer": key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "view book metadata"),
		),
		"CategoryList": key.NewBinding(
			key.WithKeys("c", "tab"),
			key.WithHelp("c", "Browse Categories"),
		),
		"ToggleItem": key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "select item"),
		),
		"PrevMenu": key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "prev menu"),
		),
		"ExitScreen": key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "exit screen"),
		),
		"DeselectAll": key.NewBinding(
			key.WithKeys("V"),
			key.WithHelp("V", "deselect all"),
		),
		"SelectAll": key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "select all"),
		),
		"ChangeLibrary": key.NewBinding(
			key.WithKeys("L"),
			key.WithHelp("L", "change library"),
		),
		"FullScreen": key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "full screen"),
		),
		"SortBy": key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "sort options"),
		),
		"EditField": key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit meta"),
		),
	}
}

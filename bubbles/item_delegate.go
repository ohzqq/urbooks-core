package bubbles

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/truncate"
)

type item struct {
	title    string
	id       string
	marked   string
	unmarked string
	selected bool
}

func (i item) Title() string       { return i.title }
func (i item) IsSelected() bool    { return i.selected }
func (i item) ID() string          { return i.id }
func (i item) FilterValue() string { return i.title }

//func (i *item) ToggleSelected()    { i.selected = !i.selected }

func (i *item) ToggleSelected() {
	i.selected = !i.selected
	if i.selected {
		i.title = strings.ReplaceAll(i.title, i.unmarked, i.marked)
	} else {
		i.title = strings.ReplaceAll(i.title, i.marked, i.unmarked)
	}
}

func (i item) Marker() string {
	if i.IsSelected() {
		return i.marked
	} else {
		return i.unmarked
	}
}

//func newItemDelegate(keys *KeyMap) list.DefaultDelegate {
//  d := list.NewDefaultDelegate()

//  d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
//    var i item
//    var title string

//    if item, ok := m.SelectedItem().(item); ok {
//      i = item
//    } else {
//      return nil
//    }
//    title = i.Title()

//    switch msg := msg.(type) {
//    case tea.KeyMsg:
//      switch {
//      case key.Matches(msg, d.keys["ToggleItem"]):
//        i.ToggleSelected()
//        m.SetItem(m.Index(), i)
//      }
//    }

//    return nil
//  }

//  //help := []key.Binding{keys.choose, keys.remove, keys.sel}

//  //d.ShortHelpFunc = func() []key.Binding {
//  //  return help
//  //}

//  //d.FullHelpFunc = func() [][]key.Binding {
//  //  return [][]key.Binding{help}
//  //}

//  d.ShowDescription = false

//  return d
//}

type itemDelegate struct {
	MultiSelect bool
	keys        KeyMap
	styles      ItemStyle
}

func (d itemDelegate) Height() int {
	return 1
}

func (d itemDelegate) Spacing() int {
	return 0
}

func (d itemDelegate) ShortHelp() []key.Binding {
	return d.keys.ShortHelp()
}

func (d itemDelegate) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{d.keys["Enter"], d.keys["FullScreen"]},
	}
}

func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.keys["ToggleItem"]):
			if d.MultiSelect {
				switch i := m.SelectedItem().(type) {
				case item:
					i.ToggleSelected()
					m.SetItem(m.Index(), i)
				}
			}
		}
	}
	return nil
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var (
		title    string
		selected bool
		style    = &d.styles
	)

	switch i := listItem.(type) {
	case item:
		title = i.Title()
		selected = i.IsSelected()
	}

	if m.Width() > 0 {
		textwidth := uint(m.Width() - style.CurrentItem.GetPaddingLeft() - style.CurrentItem.GetPaddingRight())
		title = padding.String(truncate.StringWithTail(title, textwidth, ellipsis), textwidth)
	}

	var (
		isCurrent  = index == m.Index()
		isSelected = selected
		check      = "[x] "
		uncheck    = "[ ] "
		mark       = " >> "
	)

	fn := style.NormalItem.Render

	switch d.MultiSelect {
	case true:
		mark = uncheck
		if isSelected {
			mark = check
		}
		if isCurrent {
			fn = func(s string) string {
				return style.CurrentItem.Render(mark + s)
			}
		} else if isSelected {
			fn = func(s string) string {
				return style.SelectedItem.Render(mark + s)
			}
		} else {
			fn = func(s string) string {
				return style.NormalItem.Render(mark + s)
			}
		}
	case false:
		if isCurrent {
			fn = func(s string) string {
				return style.CurrentItem.Render(s)
			}
		} else if isSelected {
			fn = func(s string) string {
				return style.SelectedItem.Render(s)
			}
		} else {
			fn = func(s string) string {
				return style.NormalItem.Render(s)
			}
		}
	}

	fmt.Fprintf(w, fn(title))
}

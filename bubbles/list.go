package bubbles

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

type listModel struct {
	list     list.Model
	items    []list.Item
	title    string
	multi    bool
	width    int
	height   int
	showHelp bool
}

func NewList() *listModel {
	return &listModel{}
}

func (l *listModel) GetSelected() []string {
	var sel []string
	for _, i := range l.list.Items() {
		switch li := i.(type) {
		case item:
			if li.IsSelected() {
				sel = append(sel, li.ID())
			}
		}
	}
	return sel
}

//func (m model) getSelectedItems() string {
//  var sel []string
//  for _, i := range m.list.Items() {
//    if item, ok := i.(item); ok {
//      if item.IsSelected() {
//        sel = append(sel, strings.TrimPrefix(item.Title(), item.Marker()))
//      }
//    }
//  }
//  return strings.Join(sel, ", ")
//}

func (l *listModel) SetHeight(h int) *listModel {
	l.height = h
	return l
}

func (l *listModel) SetWidth(w int) *listModel {
	l.width = w
	return l
}

func (l *listModel) SetTitle(t string) *listModel {
	l.title = t
	return l
}

func (l *listModel) Multi() *listModel {
	l.multi = true
	return l
}

func (l *listModel) ShowHelp() *listModel {
	l.showHelp = true
	return l
}

func (l *listModel) AppendItem(i list.Item) *listModel {
	l.items = append(l.items, i)
	return l
}

func (l *listModel) Model() list.Model {
	m := list.New(l.items,
		itemDelegate{
			MultiSelect: l.multi,
			styles:      ItemStyles(),
			keys:        DefaultKeyMap(),
		},
		l.width,
		l.height,
	)
	m.Title = l.title
	m.Styles = ListStyles()
	//m.KeyMap = listKeyMap()
	m.SetShowStatusBar(false)
	m.SetShowHelp(l.showHelp)

	return m
}

func listKeyMap() list.KeyMap {
	km := list.DefaultKeyMap()
	km.NextPage = key.NewBinding(
		key.WithKeys("right", "l", "pgdown"),
		key.WithHelp("l/pgdn", "next page"),
	)
	km.Quit = key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c", "quit"),
	)
	return km
}

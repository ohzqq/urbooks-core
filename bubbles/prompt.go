package bubbles

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Prompt struct {
	model  list.Model
	list   *listModel
	data   map[string]string
	keyMap KeyMap
	Choice string
}

func NewPrompt(title string, items map[string]string) *Prompt {
	l := NewList().
		SetTitle(title).
		SetHeight(TermHeight() - 2).
		SetWidth(TermWidth()).
		ShowHelp()

	for key, _ := range items {
		l.AppendItem(item{title: key, id: key, selected: false})
	}

	return &Prompt{
		model:  l.Model(),
		list:   l,
		data:   items,
		keyMap: DefaultKeyMap(),
	}
}

func (m *Prompt) Choose() string {
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	return m.Choice
}

func (m *Prompt) Init() tea.Cmd {
	return nil
}

func (m *Prompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap["Enter"]):
			cmds = append(cmds, m.selectListItem())
		case key.Matches(msg, m.keyMap["ExitScreen"]):
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, m.keyMap["Quit"]):
			cmds = append(cmds, tea.Quit)
		}
	case tea.WindowSizeMsg:
	case selectedListItemMsg:
		m.Choice = m.data[string(msg)]
		cmds = append(cmds, tea.Quit)
	}

	m.model, cmd = m.model.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Prompt) View() string {
	return m.model.View()
}

type selectedListItemMsg string

func (m *Prompt) selectListItem() tea.Cmd {
	return func() tea.Msg {
		var msg selectedListItemMsg
		if i, ok := m.model.SelectedItem().(item); ok {
			msg = selectedListItemMsg(i.ID())
		}
		return msg
	}
}

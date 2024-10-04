package console

import tea "github.com/charmbracelet/bubbletea"

type MainModel struct {
	current tea.Model
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return (m.current).Update(msg)
}

func (m MainModel) View() string {
	return m.current.View()
}

func (m MainModel) Init() tea.Cmd {
	return m.current.Init()
}

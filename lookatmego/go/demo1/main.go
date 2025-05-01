package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	vp    viewport.Model
	input []string
}

func (m *model) Init() tea.Cmd {
	panic("unimplemented")
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch ev := msg.(type) {
	case tea.WindowSizeMsg:
		m.vp.Height = ev.Height - 1
		m.vp.Width = ev.Width

		return m, nil
	case tea.KeyMsg:
		return m, nil
	}

	vp, cmd := m.vp.Update(msg)
	m.vp = vp

	return m, cmd
}

func (m *model) View() string {
	panic("unimplemented")
}

func (m *model) footer() string {
	return fmt.Sprintf("w: %d, h: %d", m.vp.Width, m.vp.Height)
}

func main() {
	program := tea.NewProgram(&model{})

	program.Run()
}

var (
	footerStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#5555ee"))

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#abc123"))
)

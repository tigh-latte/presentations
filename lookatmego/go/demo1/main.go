package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	h, w  int
	vp    viewport.Model
	input []string
}

// Init implements tea.Model.
func (m *model) Init() tea.Cmd {
	m.input = make([]string, 0)
	return nil
}

type buildInputMsg struct{}

func buildInput() tea.Msg {
	return buildInputMsg{}
}

// Update implements tea.Model.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch ms := msg.(type) {
	case buildInputMsg:
		content := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("2")).Render(strings.Join(m.input, ""))
		content = lipgloss.JoinHorizontal(lipgloss.Left, "so far we have typed ", content)
		m.vp.SetContent(content)
		return m, nil

	case tea.WindowSizeMsg:
		m.h = ms.Height
		m.w = ms.Width

		m.vp.Height = ms.Height - 1
		m.vp.Width = ms.Width

		return m, buildInput
	case tea.KeyMsg:
		switch ms.String() {
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
			return m, buildInput
		case "ctrl+d", "ctrl+q", "ctrl+c":
			return m, tea.Quit
		}

		m.input = append(m.input, ms.String())
		return m, buildInput
	}

	vp, cmd := m.vp.Update(msg)
	m.vp = vp

	return m, cmd
}

// View implements tea.Model.
func (m *model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.vp.View(),
		m.footer(),
	)
}

func (m *model) footer() string {
	return lipgloss.NewStyle().Width(m.w).Background(lipgloss.Color("67")).Foreground(lipgloss.Color("0")).Render(fmt.Sprintf("h: %d, w: %d", m.h, m.w))
}

func main() {
	program := tea.NewProgram(&model{})

	program.Run()
}

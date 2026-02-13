package teaui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type step int

const (
	stepType step = iota
	stepScope
	stepTitle
	stepBody
	stepConfirm
	stepDone
)

type Model struct {
	step step

	choices  []string
	filtered []string
	cursor   int

	commitType string
	scope      string
	title      string
	body       []string
	bodyPrefix string
	typeFilter textinput.Model
	input      textinput.Model
	textarea   textarea.Model

	result string

	termWidth  int
	termHeight int
}

func New(bodyPrefix string) Model {
	ti := textinput.New()
	ti.CharLimit = 72
	ti.Width = 50

	filter := textinput.New()
	filter.Placeholder = "Filter type..."
	filter.Width = 20
	filter.Focus()

	ta := textarea.New()
	ta.Placeholder = "Commit body (Ctrl+S to finish)"
	ta.ShowLineNumbers = false
	ta.SetWidth(60)
	ta.SetHeight(10)

	choices := []string{"feat", "fix", "chore", "docs"}

	return Model{
		step:       stepType,
		choices:    choices,
		filtered:   choices,
		typeFilter: filter,
		input:      ti,
		textarea:   ta,
		bodyPrefix: bodyPrefix,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) applyFilter() {
	q := strings.ToLower(m.typeFilter.Value())
	m.filtered = nil
	for _, c := range m.choices {
		if strings.Contains(c, q) {
			m.filtered = append(m.filtered, c)
		}
	}

	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		// dynamic resizing
		m.input.Width = msg.Width - 4
		m.typeFilter.Width = msg.Width - 4
		m.textarea.SetWidth(msg.Width - 4)
		m.textarea.SetHeight(msg.Height - 8)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	switch m.step {
	case stepType:
		return m.updateType(msg)
	case stepScope:
		return m.updateScope(msg)
	case stepTitle:
		return m.updateTitle(msg)
	case stepBody:
		return m.updateBody(msg)
	case stepConfirm:
		return m.updateConfirm(msg)
	}

	return m, nil
}

func (m Model) updateType(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.typeFilter, cmd = m.typeFilter.Update(msg)
	m.applyFilter()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.filtered) > 0 {
				m.commitType = m.filtered[m.cursor]
				m.input.Placeholder = "Scope (optional)"
				m.input.Focus()
				m.step = stepScope
			}
		}
	}

	return m, cmd
}

func (m Model) updateScope(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		m.scope = m.input.Value()
		m.input.Reset()
		m.input.Placeholder = "Short commit title"
		m.step = stepTitle
	}

	return m, cmd
}

func (m Model) updateTitle(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		m.title = m.input.Value()
		m.textarea.Focus()
		m.step = stepBody
	}

	return m, cmd
}

func (m Model) updateBody(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)

	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+s":
			// Recompute body with prefix
			m.body = nil
			lines := strings.Split(m.textarea.Value(), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					m.body = append(m.body, m.bodyPrefix+" "+line)
				}
			}
			m.result = m.buildCommit()
			m.step = stepConfirm
		}
	}

	return m, cmd
}

func (m Model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			m.step = stepDone
			return m, tea.Quit
		case "e":
			m.step = stepBody
			m.textarea.Focus()
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.step {
	case stepType:
		s := "Select commit type (type to filter, Enter to select)\n\n"
		s += m.typeFilter.View() + "\n\n"

		for i, c := range m.filtered {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, c)
		}
		if len(m.filtered) == 0 {
			s += "\n(no matches)"
		}
		return s

	case stepScope:
		return "Scope\n\n" + m.input.View()

	case stepTitle:
		return "Title\n\n" + m.input.View()

	case stepBody:
		return "Body (Enter = newline, Ctrl+S to finish)\n\n" + m.textarea.View()

	case stepConfirm:
		return fmt.Sprintf(
			"Commit Preview\n\n%s\n[y] submit  [e] edit body  [q] quit\n",
			m.result,
		)
	}

	return ""
}

// buildCommit ensures no leading spaces before commit type
func (m Model) buildCommit() string {
	header := m.commitType
	if m.scope != "" {
		header += "(" + m.scope + ")"
	}
	header += ": " + m.title

	out := header
	if len(m.body) > 0 {
		out += "\n\n" + strings.Join(m.body, "\n")
	}
	return out
}

// Result returns the final commit string
func (m Model) Result() string {
	return m.result
}

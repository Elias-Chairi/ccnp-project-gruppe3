package view

import (
	"bufio"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TerminalView replaces the Fyne GUI with a simple terminal interface.
type TerminalView struct {
	reader     *bufio.Reader
}

// NewTerminal creates and returns a new terminal view.
func NewTerminal() *TerminalView {
	v := &TerminalView{
		reader: bufio.NewReader(os.Stdin),
	}
	return v
}

func (v *TerminalView) Start() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error running TUI:", err)
		os.Exit(1)
	}
}

// ChangeLabel would be used by the controller (compatibility placeholder).
func (v *TerminalView) ChangeLabel(text string) {
	fmt.Println("INFO:", text)
}

// --------------------- Bubble Tea Model ---------------------

type model struct {
	message  string
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func initialModel() model {
	return model{
		message: "Welcome to the Farm Control Panel!\nPress 'q' to quit.",

		choices:  []string{"Node 1", "Node 2", "Node 3"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// Header
	s := "The DashBoard\n"

	// Iterate over choices
	for i, choice := range m.choices {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// Footer
	s += "\nPress q to quit.\n"

	// Apply style to message and add it to the output
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#0d4928ff"))

	s += "\n" + style.Render(m.message)

	return s
}

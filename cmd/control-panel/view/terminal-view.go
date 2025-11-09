package view

import (
	"fmt"
	"os"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TerminalView replaces the Fyne GUI with a simple terminal interface.
type TerminalView struct {
	Nodes []entity.Node
}

func (v *TerminalView) Start() {
	p := tea.NewProgram(initialModel(v.Nodes))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running TUI:", err)
		os.Exit(1)
	}
}

// --------------------- Bubble Tea Model ---------------------

type viewState int

const (
	mainMenu viewState = iota
	greenhouseList
	connectedMenu
	sensorView
	actuatorView
)

type model struct {
	viewState    viewState
	message      string
	choices      []string
	cursor       int
	selected     map[int]struct{}
	nodes        *[]entity.Node
	selectedNode int
}

var modelStack util.Stack[model]

func initialModel(nodes []entity.Node) tea.Model {
	return model{
		viewState:    mainMenu,
		message:      "Welcome to the Farm Control Panel!\nPress 'q' to quit.",
		choices:      []string{"Manage Greenhouses", "Exit"},
		selected:     make(map[int]struct{}),
		nodes:        &nodes,
		selectedNode: 0,
	}
}

func (m model) Init() tea.Cmd {
	modelStack.Push(m)
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

		// The "up" and "w" keys move the cursor up
		case "up", "w", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "s" keys move the cursor down
		case "down", "s", "j":
			limit := 0
			switch m.viewState {
			case actuatorView:
				if m.selectedNode < len(*m.nodes) {
					node := (*m.nodes)[m.selectedNode]
					limit = len(node.Actuators)
				}
			default:
				limit = len(m.choices)
			}
			if m.cursor < limit-1 {
				m.cursor++
			}

		case "left", "backspace", "a", "h":
			if m.viewState != mainMenu {
				nm, _ := modelStack.Pop()
				m = *nm
			}

		case "enter", "space":
			switch m.viewState {
			case actuatorView:
				// Toggle ON/OFF for the selected actuator
				if len(*m.nodes) > m.selectedNode {
					acts := (*m.nodes)[m.selectedNode].Actuators
					if len(acts) > 0 && m.cursor < len(acts) {
						act := acts[m.cursor]
						switch v := act.State.(type) {
						case bool:
							act.State = !v
							m.message = fmt.Sprintf("Toggled %s to %v", act.Type, act.State)

						case string:
							if v == "ON" {
								act.State = "OFF"
							} else {
								act.State = "ON"
							}
							m.message = fmt.Sprintf("Toggled %s to %s", act.Type, act.State)

						default:
							m.message = fmt.Sprintf("Actuator %s has unsupported state type: %T", act.Type, act.State)
						}

					}
				} else {
					m.message = fmt.Sprintf("Greenhouse %c does not exist.", 'A'+m.selectedNode)
				}
			}

		case "right", "d", "l":
			if m.viewState != actuatorView && m.viewState != sensorView {
				modelStack.Push(m)
			}
			switch m.viewState {
			case mainMenu:
				switch m.cursor {
				case 0:
					// Go to greenhouse list

					m.viewState = greenhouseList
					m.choices = make([]string, len(*m.nodes))
					for i := range *m.nodes {
						m.choices[i] = fmt.Sprintf("Greenhouse %c", 'A'+i)
					}
					m.cursor = 0
					m.message = "Viewing list of available greenhouses."
				default:
					return m, tea.Quit
				}

			case greenhouseList:
				if m.cursor < len(*m.nodes) {
					m.selectedNode = m.cursor
					m.viewState = connectedMenu
					m.choices = []string{"View Sensor Data", "View/Change Actuator Status"}
					m.cursor = 0
					m.message = fmt.Sprintf("Connected to Greenhouse %c (Demo Mode)", 'A'+m.selectedNode)
				}

			case connectedMenu:
				switch m.cursor {
				case 0:
					m.viewState = sensorView
					m.message = "Displaying sensor data"
				case 1:
					m.viewState = actuatorView
					m.cursor = 0
					m.message = "Displaying actuator status"
				}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

// This function is responsible for rendering the different views based on the current state.
func (m model) View() string {

	switch m.viewState {
	case mainMenu:
		return renderMenu(" SMART GREENHOUSE CLIENT", m.choices, m.cursor, m.message)
	case greenhouseList:
		return renderMenu("Available Greenhouses:", m.choices, m.cursor, m.message)
	case connectedMenu:
		return renderMenu("Connected to: Greenhouse", m.choices, m.cursor, m.message)
	case sensorView:
		return renderSensorView(*m.nodes, m.selectedNode)
	case actuatorView:
		return renderActuatorView(m)
	default:
		return "Unknown state"
	}
}

// Renders a menu with the given title, choices, cursor position, and message.
func renderMenu(title string, choices []string, cursor int, message string) string {
	s := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00aa55")).Render(title)
	s += "\n-----------------------------\n"

	for i, choice := range choices {
		prefix := " "
		if i == cursor {
			prefix = ">"
		}

		s += fmt.Sprintf("%s %s\n", prefix, choice)
	}

	s += "\nPress ↑/↓ and Enter to select. Press backspace to go back or q to quit.\n"
	if message != "" {
		s += "\n" + message
	}
	return s
}

func renderSensorView(nodes []entity.Node, selectedNode int) string {

	s := fmt.Sprintf(" Sensor Data for Greenhouse %c\n", 'A'+selectedNode)
	s += "--------------------------------\n"

	node := nodes[selectedNode]
	if len(node.Sensors) == 0 {
		s += "No sensors found.\n"
	} else {
		for _, sensor := range node.Sensors {
			switch v := sensor.Value.(type) {
			case float64, float32:
				s += fmt.Sprintf("[%-12s]  %.2f %s\n", sensor.Type, v, sensor.Unit)
			case int, int32, int64:
				s += fmt.Sprintf("[%-12s]  %d %s\n", sensor.Type, v, sensor.Unit)
			default:
				s += fmt.Sprintf("[%-12s]  %v %s\n", sensor.Type, v, sensor.Unit)
			}
		}
	}

	s += "\nPress q to quit or backspace to go back.\n"

	return s
}

func renderActuatorView(m model) string {
	if m.selectedNode >= len(*m.nodes) {
		return fmt.Sprintf("Greenhouse %c not found.\n", 'A'+m.selectedNode)
	}

	s := fmt.Sprintf("  Actuator Status for Greenhouse %c\n", 'A'+m.selectedNode)
	s += "-----------------------------------\n"

	node := (*m.nodes)[m.selectedNode]
	if len(node.Actuators) == 0 {
		s += "No actuators found.\n"
		return s
	}

	highlight := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#e1ce40ff"))

	for i, act := range (*m.nodes)[m.selectedNode].Actuators {

		var style lipgloss.Style

		switch v := act.State.(type) {
		case string:
			switch v {
			case "OFF":
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("#e52e19ff"))
			case "ON":
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("#5de140ff"))
			default:
				style = lipgloss.NewStyle()
			}
		default:
			style = lipgloss.NewStyle()
		}

		line := fmt.Sprintf("[%-10s]  %-3v", act.Type, act.State)
		coloredLine := style.Render(line)

		if i == m.cursor {
			s += highlight.Render(coloredLine) + "\n"
		} else {
			s += coloredLine + "\n"
		}
	}

	s += "\nUse ↑/↓ to navigate, → to toggle, ← to go back, q to quit.\n"
	s += "\n" + m.message
	return s
}

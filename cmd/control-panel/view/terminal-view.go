package view

import (
	"fmt"
	"net"
	"os"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TerminalView replaces the Fyne GUI with a simple terminal interface.
type TerminalView struct {
	nodes      []entity.Node
	teaProgram *tea.Program
}

func (t *TerminalView) AddNodes(nodes []entity.Node) {
	t.nodes = append(t.nodes, nodes...)

}

func (t *TerminalView) Start() {
	t.teaProgram = tea.NewProgram(initialModel(t.nodes))
	if _, err := t.teaProgram.Run(); err != nil {
		fmt.Println("Error running TUI:", err)
		os.Exit(1)
	}
}

type setErr error

func (t *TerminalView) FailedToConnectToServer(IP net.IP) {
	t.teaProgram.Send(setErr(fmt.Errorf("Failed to connect to server with IP %v", IP.String())))
}

func (t *TerminalView) FailedToRegisterToServer(IP net.IP) {
	t.teaProgram.Send(setErr(fmt.Errorf("Failed to register to server with IP %v", IP.String())))
}

type setLoadingMessage string

func (t *TerminalView) StartLoadingInitialNodes() {
	t.teaProgram.Send(setLoadingMessage("Loading Nodes"))
}

func (t *TerminalView) EndLoading() {
	t.teaProgram.Send(setLoadingMessage(""))
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
	nodes        *[]entity.Node
	selectedNode int
	stack        *util.Stack[model] // used to manage navigation history
	err          error
	// spinner 	 spinner.Model
	loadingmsg string
	loading    bool
}

func initialModel(nodes []entity.Node) tea.Model {
	return model{
		viewState:    mainMenu,
		message:      "Welcome to the Farm Control Panel!\nPress 'q' to quit.",
		choices:      []string{"Manage Greenhouses", "Exit"},
		nodes:        &nodes,
		stack:        &util.Stack[model]{},
		selectedNode: 0,
		loadingmsg: "",
		err: nil,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case setErr:
		m.err = error(msg)

	case setLoadingMessage:
		m.loadingmsg = string(msg)
	
	// Is it a key press?
	case tea.KeyMsg:

		switch msg.String() {	
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		if (m.err != nil || m.loadingmsg != "") {
			break
		}

		// Cool, what was the actual key pressed?
		switch msg.String() {

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
				nm, _ := m.stack.Pop()
				m = *nm
			}

		case "enter", "space":
			switch m.viewState {
			case actuatorView:
				// Makes sure the selected node actually exists
				if len(*m.nodes) > m.selectedNode {

					// Get a pointer to the selected greenhouse node.
					// Pointer to change the actuators directly
					node := &(*m.nodes)[m.selectedNode]

					// Checks that cursor is within the actuator list range.
					if m.cursor < len(node.Actuators) {
						// Pointer to the selected actuator
						act := &node.Actuators[m.cursor]
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
				m.stack.Push(m)
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
	if m.err != nil {
		return renderError(m.err)
	}

	if m.loadingmsg != "" {
		return renderLoading(m.loadingmsg)
	}

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

func renderLoading(msg string) string {
	title := "Loading"
	s := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#008eaaff")).Render(title)
	s += "\n-----------------------------\n"
	s += msg + "\n"

	s += "\nPress ↑/↓ and Enter to select. Press backspace to go back or q to quit.\n"
	return s
}

func renderError(err error) string {
	title := "Error occured"
	s := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#aa1100ff")).Render(title)
	s += "\n-----------------------------\n"
	s += fmt.Sprintf("Error: %s\n", err.Error())

	s += "\nPress ↑/↓ and Enter to select. Press backspace to go back or q to quit.\n"
	return s
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

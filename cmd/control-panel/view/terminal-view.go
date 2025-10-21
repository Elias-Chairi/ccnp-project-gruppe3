package view

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TerminalView replaces the Fyne GUI with a simple terminal interface.
type TerminalView struct {
	reader     *bufio.Reader
	nodes  []*entity.Node
}

// NewTerminal creates and returns a new terminal view.
func NewTerminalWithData(nodes []*entity.Node) *TerminalView {
	return &TerminalView{
		reader: bufio.NewReader(os.Stdin),
		nodes:  nodes,
	}
}

func (v *TerminalView) Start() {
	p := tea.NewProgram(initialModel(v.nodes))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running TUI:", err)
		os.Exit(1)
	}
}

// ChangeLabel would be used by the controller (compatibility placeholder).
func (v *TerminalView) ChangeLabel(text string) {
	fmt.Println("INFO:", text)
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
    state    viewState
    message  string
    choices  []string
    cursor   int
    selected map[int]struct{}
	nodes []*entity.Node
	selectedNode int
}


func initialModel(nodes []*entity.Node) model {
    return model{
        state:    mainMenu,
        message:  "Welcome to the Farm Control Panel!\nPress 'q' to quit.",
        choices:  []string{"Manage Greenhouses", "Exit"},
        selected: make(map[int]struct{}),
        nodes:    nodes,
		selectedNode: 0,
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

		// The "up" and "w" keys move the cursor up
		case "up", "w", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "s" keys move the cursor down
		case "down", "s", "j":
			limit := 0
			switch m.state {
			case actuatorView:
				if len(m.nodes) > 0 {
					limit = len(m.nodes[0].Actuators)
				}
			default:
				limit = len(m.choices)
			}
			if m.cursor < limit-1 {
				m.cursor++
			}


		case "left", "backspace", "a", "h": // ← Backspace go back
			switch m.state {
			case greenhouseList:
				m.state = mainMenu
				m.choices = []string{"Manage Greenhouses", "Exit"}
				m.cursor = 0
				m.message = "Returned to main menu."
			case connectedMenu:
				m.state = greenhouseList
				m.choices = []string{
					"Greenhouse A", 
					"Greenhouse B", 
					"Greenhouse C", 
				}
				m.cursor = 0
				m.message = "Returned to greenhouse list."
			case sensorView, actuatorView:
				m.state = connectedMenu
				m.choices = []string{
					"View Sensor Data",
					"View/Change Actuator Status",
				}
				m.cursor = 0
				m.message = "Returned to connected menu."
			}	

		case "right", "enter", "d", "l":
			switch m.state {
			case mainMenu:
				switch m.cursor {
				case 0:
					// Go to greenhouse list
					m.state = greenhouseList
					m.choices = []string{
						"Greenhouse A ",
						"Greenhouse B ",
						"Greenhouse C ",
					}
					m.cursor = 0

				default:
					return m, tea.Quit
				}


			case greenhouseList:
				switch m.cursor {
				case 0, 1, 2:
					
					if m.cursor >= len(m.nodes) {
						m.message = fmt.Sprintf("greenhouse %c not available.", 'A'+m.cursor)
						return m, nil
					}
				
					m.state = connectedMenu
					m.selectedNode = m.cursor
					m.choices = []string{
						"View Sensor Data",
						"View/Change Actuator Status",
					}
					m.cursor = 0
					m.message = fmt.Sprintf("Connected to: Greenhouse %c (Demo Mode)", 'A'+m.selectedNode)
				}

			case connectedMenu:
				switch m.cursor {
				case 0:
					m.state = sensorView
					m.message = "Displaying sensor data"
				case 1:
					m.state = actuatorView
					m.cursor = 0
					m.message = "Displaying actuator status"
				}

			case actuatorView:
				// Toggle ON/OFF for the selected actuator
				if len(m.nodes) > m.selectedNode {
					acts := m.nodes[m.selectedNode].Actuators
					if len(acts) > 0 && m.cursor < len(acts) {
						act := acts[m.cursor]
						if state, ok := act.State.(string); ok {
							if state == "ON" {
								act.State = "OFF"
							} else {
								act.State = "ON"
							}
							m.message = fmt.Sprintf("Toggled %s to %s", act.Type, act.State)
						}
					}	
				}	else {
					m.message = fmt.Sprintf("Greenhouse %c does not exist.", 'A'+m.selectedNode)
				}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0d4928ff"))

	switch m.state {

	case mainMenu:
		return renderMenu(" SMART GREENHOUSE CLIENT", m.choices, m.cursor, m.message, m.nodes)

	case greenhouseList:
		return renderMenu("Available Greenhouses:", m.choices, m.cursor, m.message, m.nodes)

	case connectedMenu:
		return renderMenu("Connected to: Greenhouse", m.choices, m.cursor, m.message, m.nodes)

	case sensorView:
    	return renderSensorView(m.nodes)
	case actuatorView:
    	return renderActuatorView(m)


	default:
		return style.Render("Unknown state")
	}
}

func renderMenu(title string, choices []string, cursor int, message string, nodes []*entity.Node) string {
	s := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00aa55")).Render(title)
	s += "\n-----------------------------\n"

	for i, choice := range choices {
		prefix := " "
		if i == cursor {
			prefix = ">"
		}

		if title == "Available Greenhouses:" {
			if  i >= len(nodes) {
				choice = lipgloss.NewStyle().Faint(true).Render(choice + "(unavalable)")
			}
		}
		s += fmt.Sprintf("%s %s\n", prefix, choice)
	}

	s += "\nPress ↑/↓ and Enter to select. Press backspace to go back or q to quit.\n"
	if message != "" {
		s += "\n" + message
	}
	return s
}


func renderSensorView(nodes []*entity.Node) string {
    s := " Sensor Data for Greenhouse A\n"
    s += "--------------------------------\n"

    for _, sensor := range nodes[0].Sensors {
		switch v := sensor.Value.(type) {
		case float64, float32:
			s += fmt.Sprintf("[%-12s]  %.2f %s\n", sensor.Type, v, sensor.Unit)
		case int, int32, int64:
			s += fmt.Sprintf("[%-12s]  %d %s\n", sensor.Type, v, sensor.Unit)
		default:
			s += fmt.Sprintf("[%-12s]  %v %s\n", sensor.Type, v, sensor.Unit)
		}
	}
    s += "\nPress q to quit or backspace go back."
    return s
}

func renderActuatorView(m model) string {
    s := "  Actuator Status for Greenhouse A\n"
    s += "-----------------------------------\n"

	if m.selectedNode >= len(m.nodes) {
		s += fmt.Sprintf("Greenhouse %c not found \n", 'A'+m.selectedNode)
		return s
	}

    if len(m.nodes) == m.selectedNode || len(m.nodes[m.selectedNode].Actuators) == 0 {
        s += "No actuators found.\n"
        return s
    }

    highlight := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#e1ce40ff"))

    for i, act := range m.nodes[m.selectedNode].Actuators {

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




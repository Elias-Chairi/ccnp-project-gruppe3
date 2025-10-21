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
}


func initialModel(nodes []*entity.Node) model {
    return model{
        state:    mainMenu,
        message:  "Welcome to the Farm Control Panel!\nPress 'q' to quit.",
        choices:  []string{"Manage Greenhouses", "Exit"},
        selected: make(map[int]struct{}),
        nodes:    nodes,
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
				m.choices = []string{"Greenhouse A", "Greenhouse B", "Greenhouse C"}
				m.cursor = 0
				m.message = "Returned to greenhouse list."
			case sensorView, actuatorView:
				m.state = connectedMenu
				m.choices = []string{
					"View Sensor Data",
					"View/Change Actuator Status",
					"Disconnect",
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
						"Back to Main Menu",
					}
					m.cursor = 0

				default:
					return m, tea.Quit
				}


			case greenhouseList:
				switch m.cursor {
				case 0, 1, 2:
					// Pretend to connect to one of the greenhouses
					m.state = connectedMenu
					m.choices = []string{
						"View Sensor Data",
						"View/Change Actuator Status",
						"Disconnect",
					}
					m.cursor = 0
					m.message = fmt.Sprintf("Connected to: %s (Demo Mode)", m.choices[m.cursor])

				case 3:
					m.state = mainMenu
					m.choices = []string{"List Available Greenhouses", "Connect to a greenhouse", "Exit"}
					m.cursor = 0
					m.message = "Returned to main menu."
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
				case 2:
					m.state = mainMenu
					m.choices = []string{
						"List Available Greenhouses",
						"Connect to a greenhouse",
						"Exit",
					}
					m.cursor = 0
					m.message = "Disconnected from Greenhouse. Returning to main menu..."
				}

			case actuatorView:
				// Toggle ON/OFF for the selected actuator
				if len(m.nodes) > 0 {
					acts := m.nodes[0].Actuators
					if len(acts) > 0 {
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
		return renderMenu(" SMART GREENHOUSE CLIENT", m.choices, m.cursor, m.message)

	case greenhouseList:
		return renderMenu("Available Greenhouses:", m.choices, m.cursor, m.message)

	case connectedMenu:
		return renderMenu("Connected to: Greenhouse", m.choices, m.cursor, m.message)

	case sensorView:
    	return renderSensorView(m.nodes)
	case actuatorView:
    	return renderActuatorView(m)


	default:
		return style.Render("Unknown state")
	}
}

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

    if len(m.nodes) == 0 || len(m.nodes[0].Actuators) == 0 {
        s += "No actuators found.\n"
        return s
    }

    highlight := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#e1c940ff"))

    for i, act := range m.nodes[0].Actuators {
        line := fmt.Sprintf("[%-10s]  %-3v", act.Type, act.State)
        if i == m.cursor {
            s += highlight.Render(line) + "\n"
        } else {
            s += line + "\n"
        }
    }

    s += "\nUse ↑/↓ to navigate, → to toggle, ← to go back, q to quit.\n"
    s += "\n" + m.message
    return s
}




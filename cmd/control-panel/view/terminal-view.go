package view

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util"

	"github.com/charmbracelet/bubbles/spinner"
	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TerminalView replaces the Fyne GUI with a simple terminal interface.
type TerminalView struct {
	nodes      []entity.Node
	teaProgram *tea.Program
}

// Start launches the terminal-based user interface.
func (t *TerminalView) Start() {
	t.teaProgram = tea.NewProgram(initialModel())
	if _, err := t.teaProgram.Run(); err != nil {
		fmt.Println("Error running TUI:", err)
		os.Exit(1)
	}
}

type setNodes []entity.Node

// SetInitialNodes sets the initial list of nodes to be displayed in the terminal after registration.
//
// Panics if the tea program is not started.
func (t *TerminalView) SetInitialNodes(nodes []entity.Node) {
	t.nodes = nodes
	t.teaProgram.Send(setNodes(nodes))
}

type setErr error

// FailedToConnectToServer notifies the terminal view of a failed connection attempt to the server.
//
// Panics if the tea program is not started.
func (t *TerminalView) FailedToConnectToServer(IP net.IP) {
	t.teaProgram.Send(setErr(fmt.Errorf("failed to connect to server with IP %v", IP.String())))
}

// FailedToRegisterToServer notifies the terminal view of a failed registration attempt to the server.
//
// Panics if the tea program is not started.
func (t *TerminalView) FailedToRegisterToServer(IP net.IP) {
	t.teaProgram.Send(setErr(fmt.Errorf("failed to register to server with IP %v", IP.String())))
}

type setLoadingMessage string

// EndLoading clears the loading message in the terminal view.
//
// Panics if the tea program is not started.
func (t *TerminalView) EndLoading() {
	t.teaProgram.Send(setLoadingMessage(""))
}

// --------------------- Bubble Tea Model ---------------------

type viewState int

const (
	mainMenu viewState = iota
	greenhouseList
	nodeView
)

type model struct {
	viewState    viewState
	message      string
	choices      []string
	cursor       int
	nodes        []entity.Node
	selectedNode int
	stack        *util.Stack[model] // used to manage navigation history
	err          error
	loadingmsg string
	editingActuator int
	input textinput.Model
	isEditingNumber bool
	spinner spinner.Model
	pendingActuator map[int]bool
}

type actuatorResponseMsg struct {
	Index int
}

func initialModel() tea.Model {

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	return model{
		viewState:    	 mainMenu,
		message:      	 "Welcome to the Farm Control Panel!\nPress 'q' to quit.",
		choices:      	 []string{"Manage Greenhouses", "Exit"},
		nodes:        	 nil,
		stack:        	 &util.Stack[model]{},
		spinner:         sp,
		pendingActuator: map[int]bool{},
		selectedNode: 	 0,
		loadingmsg:   	 "Loading Nodes",
		err:          	 nil,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	if cmd != nil {
		return m, cmd
	}


	if m.isEditingNumber {
		var inputCmd tea.Cmd
		m.input, inputCmd = m.input.Update(msg)

		switch key := msg.(type) {
		case tea.KeyMsg:
			switch key.String() {
			case "enter":
				var num int
				fmt.Sscanf(m.input.Value(), "%d", &num)
				m.nodes[m.selectedNode].Actuators[m.editingActuator].State = num
				m.isEditingNumber = false
				
				// Start spinner wait
                m.pendingActuator[m.editingActuator] = true

                return m, tea.Batch(
                    inputCmd,
                    tea.Tick(time.Second, func(t time.Time) tea.Msg {
                        return actuatorResponseMsg{Index: m.editingActuator}
                    }),
                )

			case "esc", "escape":
				m.isEditingNumber = false
				return m, nil
			}
		}
		return m, inputCmd
	}

	//main switchmsg
	switch msg := msg.(type) {

	case actuatorResponseMsg:
    	delete(m.pendingActuator, msg.Index)
    	return m, nil
 	

	case setNodes:
		m.nodes = ([]entity.Node)(msg)

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

		if m.err != nil || m.loadingmsg != "" {
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
			case nodeView:
				if m.selectedNode < len(m.nodes) {
					limit = len(m.nodes[m.selectedNode].Actuators)
				}
			default:
				limit = len(m.choices)
			}
			if m.cursor < limit-1 {
				m.cursor++
			}

		case "left", "a", "h":
			if m.viewState != mainMenu {
				nm, _ := m.stack.Pop()
				m = *nm
			}

		case "right", "d", "l":
    		switch m.viewState {

			case mainMenu:
				// same behavior as enter
				if m.cursor == 0 {
					m.stack.Push(m)
					m.viewState = greenhouseList
					m.choices = make([]string, len(m.nodes))
					for i := range m.nodes {
						m.choices[i] = fmt.Sprintf("Greenhouse %c", 'A'+i)
					}
					m.cursor = 0
				} else {
					return m, tea.Quit
				}

			case greenhouseList:
				if m.cursor < len(m.nodes) {
					m.stack.Push(m)
					m.selectedNode = m.cursor
					m.viewState = nodeView
					m.cursor = 0
				}

			case nodeView:
				return m, nil // Disable right key in node view to prevent navigation beyond leaf view
			}	

		case "backspace":
			if m.viewState == nodeView {
				node := &m.nodes[m.selectedNode]
				if m.cursor < len(node.Actuators) {
					act := &node.Actuators[m.cursor]

					switch v := act.State.(type) {

					case int:
						if v <= 10 {
							if v > 0 {
								act.State = v - 1
							}
						} else {
							act.State = v - 10
						}

					case int32:
						iv := int(v)
						if iv <= 10 {
							if iv > 0 {
								act.State = iv - 1
							}
						} else {
							act.State = iv - 10
						}

					case int64:
						iv := int(v)
						if iv <= 10 {
							if iv > 0 {
								act.State = iv - 1
							}
						} else {
							act.State = iv - 10
						}
					}
				}
				return m, nil
			}	

		case "enter", " ":
			switch m.viewState {

			case mainMenu:
				if m.cursor == 0 {
					m.stack.Push(m)
					m.viewState = greenhouseList
					m.choices = make([]string, len(m.nodes))
					for i := range m.nodes {
						m.choices[i] = fmt.Sprintf("Greenhouse %c", 'A'+i)
					}
					m.cursor = 0
				} else {
					return m, tea.Quit
				}

			case greenhouseList:
				if m.cursor < len(m.nodes) {
					m.stack.Push(m)
					m.selectedNode = m.cursor
					m.viewState = nodeView
					m.cursor = 0
				}

			case nodeView:
				node := &m.nodes[m.selectedNode]
				if m.cursor < len(node.Actuators) {
					act := &node.Actuators[m.cursor]

					switch v := act.State.(type) {

					// BOOL actuator
					case bool:
						act.State = !v
					// Mark pending and start spinner + fake delay
					m.pendingActuator[m.cursor] = true
					return m, tea.Batch(
						m.spinner.Tick,
						tea.Tick(time.Second, func(t time.Time) tea.Msg {
							return actuatorResponseMsg{Index: m.cursor}
						}),
					)
					

                    case int, int32, int64:
                        // Start editing mode
                        m.isEditingNumber = true
                        m.editingActuator = m.cursor

                        m.input = textinput.New()
                        m.input.Placeholder = "Enter number"
                        m.input.SetValue(fmt.Sprintf("%v", act.State))
                        m.input.Focus()

                        return m, textinput.Blink
					}
				}
			}
		}
	}

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
	case nodeView:
		return renderNodeView(m)
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
	title := "Error occurred"
	s := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#aa1100ff")).Render(title)
	s += "\n-----------------------------\n"
	s += fmt.Sprintf("Error: %s\n", err.Error())

	s += "\nPress ↑/↓ and Enter to select. Press ← to go back or q to quit.\n"
	return s
}


// Renders a menu with the given title, choices, cursor position, and message.
func renderMenu(title string, choices []string, cursor int, message string) string {
	s := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00aa55")).Render(title)
	s += "\n-----------------------------\n"

	var menuHighlight = lipgloss.NewStyle().
    Bold(true).
    Background(lipgloss.Color("#000ed6ff")).
    Foreground(lipgloss.Color("#ffffffff"))
	menuHighlight = menuHighlight.Width(40)
	

	for i, choice := range choices {

		line := fmt.Sprintf("  %s", choice)

		if i == cursor {
			// highlight line
			s += menuHighlight.Render(line) + "\n"
		} else {
			s += line + "\n"
		}
	}

	s += "\nPress ↑/↓ and Enter to select. Press ← to go back or q to quit.\n"
	if message != "" {
		s += "\n" + message
	}
	return s
}

func renderNodeView(m model) string {
	node := m.nodes[m.selectedNode]

	title := fmt.Sprintf(" Greenhouse %c - Node Overview\n", 'A'+m.selectedNode)
	header := title + "----------------------------------------------\n\n"

	//  LEFT COLUMN: SENSORS

	var sensors string
	sensors += lipgloss.NewStyle().Bold(true).Underline(true).Render("Sensors") + "\n"

	if len(node.Sensors) == 0 {
		sensors += "  No sensors found.\n"
	} else {
		for _, sensor := range node.Sensors {
			var line string
			switch v := sensor.Value.(type) {
			case float64, float32:
				line = fmt.Sprintf("  %-12s : %.2f %s", sensor.Type, v, sensor.Unit)
			case int, int32, int64:
				line = fmt.Sprintf("  %-12s : %d %s", sensor.Type, v, sensor.Unit)
			default:
				line = fmt.Sprintf("  %-12s : %v %s", sensor.Type, v, sensor.Unit)
			}
			sensors += line + "\n"
		}
	}

	// pad the sensors column so both columns line up
	sensorsStyle := lipgloss.NewStyle().Width(40)
	sensors = sensorsStyle.Render(sensors)


	//  RIGHT COLUMN: ACTUATORS

	highlight := lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("#000ed6ff")).
		Foreground(lipgloss.Color("#ffffffff"))
	highlight = highlight.Width(40)

	var actuators string
	actuators += lipgloss.NewStyle().Bold(true).Underline(true).Render("Actuators") + "\n"

	for i, act := range node.Actuators {

		var (
			value string
			style lipgloss.Style
		)

		switch v := act.State.(type) {

		case bool:
			if v {
				value = "ON"
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("#5de140ff"))
			} else {
				value = "OFF"
				style = lipgloss.NewStyle().Foreground(lipgloss.Color("#e52e19ff"))
			}

		case int, int32, int64:
			num := fmt.Sprintf("%v", v)

			// If editing, show the text input
			if m.isEditingNumber && i == m.editingActuator {
				value = m.input.View()
			} else if act.Unit != "" {
				value = num + " " + act.Unit
			} else {
				if num != "0" {
					value = "ON"
				} else {
					value = "OFF"
				}
			}

		}

		spinnerStr := ""
		if m.pendingActuator[i] {
			spinnerStr = " " + m.spinner.View()
		}

		line := fmt.Sprintf("  %-12s : %s%s", act.Type, value, spinnerStr)


		// actuator-specific colors
		styled := style.Render(line)

		// highlight background without overriding colors
		if i == m.cursor {
			styled = highlight.Render(styled)
		}

		actuators += styled + "\n"
	}

	actuatorsStyle := lipgloss.NewStyle().Width(40)
	actuators = actuatorsStyle.Render(actuators)

	//  JOIN BOTH COLUMNS
	combined := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sensors,
		actuators,
	)


	footer := "\n----------------------------------------------\n" +
			" ↑/↓ navigate actuators | enter/space toggle | ← back | q quit\n"

	if m.message != "" {
		footer += "\n" + m.message + "\n"
	}

	return header + combined + footer
}



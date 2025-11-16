package view

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/control-panel/connectionhandler"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	util "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util/general"

	"github.com/charmbracelet/bubbles/spinner"
	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TerminalView replaces the Fyne GUI with a simple terminal interface.
type TerminalView struct {
	nodes       []entity.Node
	teaProgram  *tea.Program
	ConnHandler *connectionhandler.ConnectionHandler
}

// Start launches the terminal-based user interface.
func (t *TerminalView) Start() {
	t.teaProgram = tea.NewProgram(initialModel(t.ConnHandler))
	if _, err := t.teaProgram.Run(); err != nil {
		fmt.Println("Error running TUI:", err)
		os.Exit(1)
	}
}

type setNodes struct {
	nodes []entity.Node
}

// SetInitialNodes sets the initial list of nodes to be displayed in the terminal after registration.
//
// Panics if the tea program is not started.
func (t *TerminalView) SetInitialNodes(nodes []entity.Node) {
	t.nodes = nodes
	t.teaProgram.Send(setNodes{nodes: nodes})
}

type setErr struct {
	err error
}

// FailedToConnectToServer notifies the terminal view of a failed connection attempt to the server.
//
// Panics if the tea program is not started.
func (t *TerminalView) FailedToConnectToServer(IP net.IP) {
	t.teaProgram.Send(setErr{err: fmt.Errorf("failed to connect to server with IP %v", IP.String())})
}

// FailedToRegisterToServer notifies the terminal view of a failed registration attempt to the server.
//
// Panics if the tea program is not started.
func (t *TerminalView) FailedToRegisterToServer(IP net.IP) {
	t.teaProgram.Send(setErr{err: fmt.Errorf("failed to register to server with IP %v", IP.String())})
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
	mainMenuView viewState = iota
	nodeListView
	nodeInfoView
)

type model struct {
	viewState       viewState
	message         string
	choices         []string
	cursor          int
	nodes           []entity.Node
	selectedNode    int
	stack           *util.Stack[model] // used to manage navigation history
	err             error
	loadingmsg      string
	editingActuator int
	input           textinput.Model
	isEditingNumber bool
	spinner         spinner.Model
	pendingActuator map[int]bool
	connHandler     *connectionhandler.ConnectionHandler
}

func (m *model) isActuatorLocked(index int) bool {
	return m.pendingActuator[index]
}

type actuatorResponseMsg struct {
	Index int
}

func initialModel(connHandler *connectionhandler.ConnectionHandler) tea.Model {

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	return model{
		viewState:       mainMenuView,
		message:         "Welcome to the Farm Control Panel!\nPress 'q' to quit.",
		choices:         []string{"Manage Greenhouses", "Exit"},
		nodes:           nil,
		stack:           &util.Stack[model]{},
		spinner:         sp,
		pendingActuator: map[int]bool{},
		selectedNode:    0,
		loadingmsg:      "Loading Nodes",
		err:             nil,
		connHandler:     connHandler,
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
				raw := m.input.Value()
				var num int

				if raw == "" {
					num = 0
				} else {
					if _, err := fmt.Sscanf(raw, "%d", &num); err != nil {
						// Invalid input → escape editing mode
						m.isEditingNumber = false
						return m, nil
					}
				}

				m.nodes[m.selectedNode].Actuators[m.editingActuator].State = num
				m.isEditingNumber = false

				// Start spinner wait
				m.pendingActuator[m.editingActuator] = true

				return m, tea.Batch(
					inputCmd,
					m.spinner.Tick,
					tea.Tick(time.Second, func(t time.Time) tea.Msg {
						return actuatorResponseMsg{Index: m.editingActuator}
					}),
				)

			case "esc", "escape":
				m.isEditingNumber = false
				return m, nil

			case "ctrl+c":
				return m, tea.Quit
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
		m.nodes = msg.nodes

	case setErr:
		m.err = msg.err

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
			case nodeInfoView:
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
			if m.viewState == nodeInfoView && m.isActuatorLocked(m.cursor) {
				return m, nil // cannot go back while this actuator is waiting
			}
			if m.viewState != mainMenuView {
				nm, _ := m.stack.Pop()
				m = *nm
			}

		case "right", "d", "l":
			switch m.viewState {

			case mainMenuView:
				// same behavior as enter
				if m.cursor == 0 {
					m.stack.Push(m)
					m.viewState = nodeListView
					m.choices = make([]string, len(m.nodes))
					for i := range m.nodes {
						m.choices[i] = fmt.Sprintf("Greenhouse %c", 'A'+i)
					}
					m.cursor = 0
				} else {
					return m, tea.Quit
				}

			case nodeListView:
				if m.cursor < len(m.nodes) {
					m.stack.Push(m)
					m.selectedNode = m.cursor
					m.viewState = nodeInfoView
					m.cursor = 0
				}

			case nodeInfoView:
				return m, nil // Disable right key in node view to prevent navigation beyond leaf view
			}

		case "backspace":
			if m.isEditingNumber {
				m.input, _ = m.input.Update(msg)
				return m, nil
			}

		case "enter", " ":
			switch m.viewState {

			case mainMenuView:
				if m.cursor == 0 {
					m.stack.Push(m)
					m.viewState = nodeListView
					m.choices = make([]string, len(m.nodes))
					for i := range m.nodes {
						m.choices[i] = fmt.Sprintf("Greenhouse %c", 'A'+i)
					}
					m.cursor = 0
				} else {
					return m, tea.Quit
				}

			case nodeListView:
				if m.cursor < len(m.nodes) {
					m.stack.Push(m)
					m.selectedNode = m.cursor
					m.viewState = nodeInfoView
					m.cursor = 0
				}

			case nodeInfoView:
				node := &m.nodes[m.selectedNode]
				if m.cursor < len(node.Actuators) {
					// Block interaction if this actuator is waiting for response
					if m.isActuatorLocked(m.cursor) {
						return m, nil
					}

					act := &node.Actuators[m.cursor]

					switch v := act.State.(type) {

					// BOOL actuator
					case bool:
						// Send command to toggle actuator state
						m.connHandler.SendCommand(m.nodes[m.selectedNode].ID, act.ID, !v)
						// act.State = !v

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

	case mainMenuView:
		return renderMenu(" SMART GREENHOUSE CLIENT", m.choices, m.cursor, m.message)
	case nodeListView:
		return renderMenu("Available Greenhouses:", m.choices, m.cursor, m.message)
	case nodeInfoView:
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

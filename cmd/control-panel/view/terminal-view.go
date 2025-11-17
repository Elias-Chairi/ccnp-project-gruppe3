package view

import (
	"fmt"
	"net"
	"os"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
	util "github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util/general"

	"github.com/charmbracelet/bubbles/spinner"
	textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Controller interface {
	SendCommand(nodeID uint8, actuatorID uint8, state any) error
}

// TerminalView replaces the Fyne GUI with a simple terminal interface.
type TerminalView struct {
	nodes      []entity.Node
	teaProgram *tea.Program
	Controller Controller
}

// Start launches the terminal-based user interface.
func (t *TerminalView) Start() {
	t.teaProgram = tea.NewProgram(initialModel(t.Controller))
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

func (t *TerminalView) SensorUpdate(nodeID uint8, sensorID uint8, value any) {
	t.teaProgram.Send(sensorUpdateMsg{
		nodeID:   nodeID,
		sensorID: sensorID,
		value:    value,
	})
}

func (t *TerminalView) ActuatorUpdate(nodeID uint8, actuatorID uint8, state any) {
	t.teaProgram.Send(actuatorUpdateMsg{
		nodeID:     nodeID,
		actuatorID: actuatorID,
		state:      state,
	})
}

func (t *TerminalView) ActuatorCommandResponse(nodeID uint8, actuatorID uint8, state any, err error) {
	if err != nil {
		return
	}
	t.teaProgram.Send(actuatorResponseMsg{
		actuatorUpdateMsg: actuatorUpdateMsg{
			nodeID:     nodeID,
			actuatorID: actuatorID,
			state:      state,
		},
	})
}

// --------------------- Bubble Tea Model ---------------------

type viewState int

const (
	mainMenuView viewState = iota
	nodeListView
	nodeInfoView
)

type model struct {
	viewState        viewState
	message          string
	choices          []string
	cursor           int
	nodes            []entity.Node
	selectedNode     int
	stack            *util.Stack[model] // used to manage navigation history
	err              error
	loadingmsg       string
	editingActuators map[int]bool
	inputFields      map[int]textinput.Model
	spinner          spinner.Model
	pendingActuator  map[int]bool
	controller       Controller
}

func (m *model) isActuatorLocked(index int) bool {
	return m.pendingActuator[index]
}

type sensorUpdateMsg struct {
	nodeID   uint8
	sensorID uint8
	value    any
}

type actuatorUpdateMsg struct {
	nodeID     uint8
	actuatorID uint8
	state      any
}

type actuatorResponseMsg struct {
	actuatorUpdateMsg
}

func (m *model) getNodeByID(id uint8) *entity.Node {
	for i := range m.nodes {
		if m.nodes[i].ID == id {
			return &m.nodes[i]
		}
	}
	return nil
}

func (m *model) updateSensor(nodeID uint8, sensorID uint8, value any) {
	node := m.getNodeByID(nodeID)
	for i := range node.Sensors {
		if node.Sensors[i].ID == sensorID {
			node.Sensors[i].Value = value
			break
		}
	}
}

func (m *model) updateActuator(nodeID uint8, actuatorID uint8, state any) int {
	node := m.getNodeByID(nodeID)
	for i := range node.Actuators {
		if node.Actuators[i].ID == actuatorID {
			node.Actuators[i].State = state
			return i
		}
	}
	return -1
}

func initialModel(controller Controller) tea.Model {

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	return model{
		viewState:        mainMenuView,
		message:          "Welcome to the Farm Control Panel!\nPress 'q' to quit.",
		choices:          []string{"Manage Greenhouses", "Exit"},
		nodes:            nil,
		stack:            &util.Stack[model]{},
		spinner:          sp,
		pendingActuator:  map[int]bool{},
		editingActuators: map[int]bool{},
		inputFields:      map[int]textinput.Model{},
		selectedNode:     0,
		loadingmsg:       "Loading Nodes",
		err:              nil,
		controller:       controller,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)

	// Handle editing for ANY actuator
	var editCmds []tea.Cmd

	for idx := range m.editingActuators {

		var cmd tea.Cmd
		m.inputFields[idx], cmd = m.inputFields[idx].Update(msg)
		if cmd != nil {
			editCmds = append(editCmds, cmd)
		}

		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {

			case "enter":
				raw := m.inputFields[idx].Value()
				original := m.nodes[m.selectedNode].Actuators[idx].State

				var state any
				switch original.(type) {
				case int, int32, int64:
					var num int
					if _, err := fmt.Sscanf(raw, "%d", &num); err != nil {
						// Invalid input, exit editing mode
						delete(m.editingActuators, idx)
						delete(m.inputFields, idx)
						break
					}
					state = num

				case float32, float64:
					var num float64
					if _, err := fmt.Sscanf(raw, "%f", &num); err != nil {
						// Invalid input, exit editing mode
						delete(m.editingActuators, idx)
						delete(m.inputFields, idx)
						break
					}
					if num < 0 {
						num = 0
					}
					if num > 1 {
						num = 1
					}
					state = num
				}

				delete(m.editingActuators, idx)
				delete(m.inputFields, idx)

				m.pendingActuator[idx] = true

				_ = m.controller.SendCommand(m.nodes[m.selectedNode].ID, m.nodes[m.selectedNode].Actuators[idx].ID, state)

			case "esc", "escape":
				delete(m.editingActuators, idx)
				delete(m.inputFields, idx)

			case "ctrl+c":
				return m, tea.Quit
			}
		}
	}

	// Return all collected edit commands AND spinner tick
	if len(editCmds) > 0 {
		return m, tea.Batch(editCmds...)
	}

	//main switchmsg
	switch msg := msg.(type) {

	case sensorUpdateMsg: // update sensor value
		m.updateSensor(msg.nodeID, msg.sensorID, msg.value)

	case actuatorUpdateMsg: // update actuator state
		m.updateActuator(msg.nodeID, msg.actuatorID, msg.state)

	case actuatorResponseMsg: // stop spinner and update actuator state
		index := m.updateActuator(msg.nodeID, msg.actuatorID, msg.state)
		delete(m.pendingActuator, index)
		return m, m.spinner.Tick

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

		// Block navigation and selection while editing ANY actuator
		if len(m.editingActuators) > 0 {
			return m, nil
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
			if m.isActuatorLocked(m.cursor) {
				return m, nil
			}
			return m, nil

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
						_ = m.controller.SendCommand(m.nodes[m.selectedNode].ID, act.ID, !v)

						// Mark pending and start spinner + fake delay
						m.pendingActuator[m.cursor] = true
						return m, m.spinner.Tick

					case int, int32, int64:
						// Start editing mode
						ti := textinput.New()
						ti.Placeholder = "Enter number"
						ti.SetValue(fmt.Sprintf("%v", act.State))
						ti.Focus()

						m.inputFields[m.cursor] = ti
						m.editingActuators[m.cursor] = true
						return m, textinput.Blink

					case float32, float64:
						// Start editing mode
						ti := textinput.New()
						ti.Placeholder = "Enter value (0–1)"
						ti.SetValue(fmt.Sprintf("%v", v))
						ti.Focus()

						m.inputFields[m.cursor] = ti
						m.editingActuators[m.cursor] = true
						return m, tea.Batch(textinput.Blink, m.spinner.Tick)

					}
				}
			}
		}
	}

	return m, tea.Batch(spinnerCmd)
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

		case float32, float64:
			var f float64
			switch x := v.(type) {
			case float32:
				f = float64(x)
			case float64:
				f = x
			}

			// While editing, show text input
			if m.editingActuators[i] {
				value = m.inputFields[i].View()
			} else {
				percent := int(f * 100)
				value = fmt.Sprintf("%d %%", percent)
			}

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
			if m.editingActuators[i] {
				value = m.inputFields[i].View()
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

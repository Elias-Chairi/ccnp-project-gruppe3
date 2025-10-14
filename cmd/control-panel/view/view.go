//go:build gui

package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/control-panel/controller"
)

type view struct {
	app        fyne.App
	window     fyne.Window
	controller controller.Controller
	label      *widget.Label
}

func New() *view {
	v := &view{}
	v.controller = controller.New(v)
	v.setupUI()
	return v
}

func (v *view) setupUI() {
	v.app = app.New()
	v.window = v.app.NewWindow("Hello World")

	v.label = widget.NewLabel("Control Panel")
	v.label.Alignment = fyne.TextAlignCenter

	leftCol := widget.NewLabel("Left Column")
	middleCol := container.NewVBox(
		widget.NewButton("Click Me", v.controller.TestButtonClicked),
	)
	rightCol := widget.NewLabel("Right Column")

	columns := container.NewGridWithColumns(3,
		leftCol,
		middleCol,
		rightCol,
	)

	content := container.NewVBox(
		v.label,
		columns,
	)

	v.window.SetContent(content)
	v.window.Resize(fyne.NewSize(400, 200))
}

func (v *view) ChangeLabel(text string) {
	v.label.SetText(text)
}

func (v *view) Start() {
	v.window.ShowAndRun()
}

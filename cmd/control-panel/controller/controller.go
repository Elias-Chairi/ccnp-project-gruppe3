package controller

type Controller interface {
	TestButtonClicked()
}

type View interface {
	ChangeLabel(text string)
}

type controller struct {
	view View
}

func New(view View) Controller {
	return &controller{view: view}
}

func (c *controller) TestButtonClicked() {
	c.view.ChangeLabel("Button Clicked!!!!!!!!!!!!!!!")
}

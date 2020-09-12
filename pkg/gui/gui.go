package gui

import (
	"github.com/jesseduffield/gocui"
)

var OverlappingEdges = false

type Gui struct {
	g *gocui.Gui
}

func NewGui() (*Gui, error) {
	return &Gui{}, nil
}

func (gui *Gui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal, OverlappingEdges)
	if err != nil {
		return err
	}
	defer g.Close()

	gui.g = g

	// shows border select
	g.Highlight = true

	// Set ColorScheme
	g.SelFgColor = gocui.ColorGreen
	g.BgColor = gocui.ColorBlack
	g.FgColor = gocui.ColorDefault

	// Allow mouse events
	g.Mouse = true
	// Set Manager
	g.SetManager(gocui.ManagerFunc(gui.layout))

	// Set Keybindings
	err = gui.SetKeybindings(g)
	if err != nil {
		return err
	}

	err = g.MainLoop()
	return err
}

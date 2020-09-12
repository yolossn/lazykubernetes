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

	g.SetManager(gocui.ManagerFunc(gui.layout))
	g.Mouse = true

	err = g.MainLoop()
	return err
}

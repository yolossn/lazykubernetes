package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/yolossn/lazykubernetes/pkg/client"
)

var OverlappingEdges = false

type Gui struct {
	g         *gocui.Gui
	k8sClient *client.K8s
}

func NewGui(k8sClient *client.K8s) (*Gui, error) {
	return &Gui{k8sClient: k8sClient}, nil
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

	// Init render
	go gui.reRenderNamespace()

	err = g.MainLoop()
	return err
}

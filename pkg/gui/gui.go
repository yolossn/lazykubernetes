package gui

import (
	"sync"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/yolossn/lazykubernetes/pkg/client"
)

var OverlappingEdges = false

type data struct {
	NamespaceData []client.NamespaceInfo
	nsMux         sync.RWMutex
	PodData       []client.PodInfo
	rsMux         sync.RWMutex
}

type Gui struct {
	g           *gocui.Gui
	k8sClient   *client.K8s
	data        *data
	panelStates *panelStates
}

func NewGui(k8sClient *client.K8s) (*Gui, error) {

	// NewData
	data := &data{}
	panelStates := NewPanelStates()
	return &Gui{k8sClient: k8sClient, data: data, panelStates: panelStates}, nil
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

	// Update Namespace data
	go gui.WatchNamespace()
	go gui.WatchPods()

	// reRender
	go gui.goEvery(time.Second, gui.reRenderNamespace)
	go gui.goEvery(time.Second, gui.reRenderResource)
	err = g.MainLoop()
	return err
}

func (gui *Gui) goEvery(interval time.Duration, function func() error) {
	// currentSessionIndex := gui.State.SessionIndex
	_ = function() // time.Tick doesn't run immediately so we'll do that here // TODO: maybe change
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			// if gui.State.SessionIndex > currentSessionIndex {
			// 	return
			// }
			_ = function()
		}
	}()
}

func (gui *Gui) FindSelectedLine(v *gocui.View, itemCount int) int {
	_, cy := v.Cursor()
	_, oy := v.Origin()

	selectedLine := cy - oy

	if selectedLine < 0 {
		return 0
	}

	if selectedLine > itemCount-1 {
		return itemCount - 1
	}
	return selectedLine
}

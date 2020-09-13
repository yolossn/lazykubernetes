package gui

import (
	"sync"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/yolossn/lazykubernetes/pkg/client"
)

var OverlappingEdges = false

type data struct {
	NamespaceData  []client.NamespaceInfo
	nsMux          sync.RWMutex
	PodData        []client.PodInfo
	JobData        []client.JobInfo
	DeploymentData []client.DeploymentInfo
	ServiceData    []client.ServiceInfo
	SecretData     []client.SecretInfo
	ConfigMapData  []client.ConfigMapInfo
	NodeData       []client.NodeInfo
	rsMux          sync.RWMutex
	nodeMux        sync.RWMutex
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
	go gui.WatchNodes()

	// reRender
	go gui.goEvery(time.Second, gui.reRenderNamespace)
	go gui.goEvery(time.Second, gui.reRenderResource)
	go gui.goEvery(time.Second, gui.reRenderClusterInfo)
	go gui.goEvery(time.Second, gui.reRenderNodeInfo)

	// highlight cluster view on start
	go gui.highlightClusterInfoView()

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

func (gui *Gui) changeSelectedLine(line *int, total int, up bool) {
	if up {
		if *line == -1 || *line == 0 {
			return
		}

		*line -= 1
	} else {
		if *line == -1 || *line == total-1 {
			return
		}

		*line += 1
	}
}

func (gui *Gui) focusPoint(selectedX int, selectedY int, lineCount int, v *gocui.View) error {
	if selectedY < 0 || selectedY > lineCount {
		return nil
	}
	ox, oy := v.Origin()
	originalOy := oy
	cx, cy := v.Cursor()
	originalCy := cy
	_, height := v.Size()

	ly := Max(height-1, 0)

	windowStart := oy
	windowEnd := oy + ly

	if selectedY < windowStart {
		oy = Max(oy-(windowStart-selectedY), 0)
	} else if selectedY > windowEnd {
		oy += (selectedY - windowEnd)
	}

	if windowEnd > lineCount-1 {
		shiftAmount := (windowEnd - (lineCount - 1))
		oy = Max(oy-shiftAmount, 0)
	}

	if originalOy != oy {
		_ = v.SetOrigin(ox, oy)
	}

	cy = selectedY - oy
	if originalCy != cy {
		_ = v.SetCursor(cx, selectedY-oy)
	}
	return nil
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

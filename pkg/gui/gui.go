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

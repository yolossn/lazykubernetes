package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"sigs.k8s.io/yaml"
)

func (gui *Gui) getNamespaceView() *gocui.View {
	v, _ := gui.g.View("namespace")
	return v
}

func (gui *Gui) onNamespaceClick(g *gocui.Gui, v *gocui.View) error {

	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	infoView := gui.getInfoView()
	infoView.Tabs = getNamespaceInfoTabs()

	// Find selectedLine
	gui.panelStates.Namespace.SelectedLine = gui.FindSelectedLine(v, len(gui.data.NamespaceData))
	fmt.Fprintln(infoView, gui.panelStates.Namespace.SelectedLine)
	err := gui.handleNSSelect(v)
	if err != nil {
		return err
	}
	return gui.reRenderResource()
}

func (gui *Gui) handleNSSelect(v *gocui.View) error {
	infoView := gui.getInfoView()
	ns := gui.getCurrentNS()
	data, err := gui.k8sClient.GetNamespace(ns)
	if err != nil {
		return err
	}

	err = gui.focusPoint(0, gui.panelStates.Namespace.SelectedLine, len(gui.data.NamespaceData), v)
	if err != nil {
		return err
	}

	infoView.Clear()
	output, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Fprintln(infoView, string(output))
	return nil
}

func (gui *Gui) getCurrentNS() string {
	var ns string

	if len(gui.data.NamespaceData) > 0 {
		if gui.panelStates.Namespace.SelectedLine > len(gui.data.NamespaceData) {
			gui.panelStates.Namespace.SelectedLine = gui.panelStates.Namespace.SelectedLine - len(gui.data.NamespaceData)
		}
		ns = gui.data.NamespaceData[gui.panelStates.Namespace.SelectedLine].Name

	}
	return ns
}

// func (gui *Gui) updateAndWatchNamespaceData() error {
func (gui *Gui) WatchNamespace() error {
	// Init fetch data
	_ = gui.updateNSData()
	// TODO:Handle error
	// Wait for namespace events and update data
	eventInterface, _ := gui.k8sClient.WatchNamespace()
	for {
		_ = <-eventInterface.ResultChan()
		_ = gui.updateNSData()
	}
	return nil
}

func (gui *Gui) updateNSData() error {
	gui.data.nsMux.Lock()
	defer gui.data.nsMux.Unlock()
	ns, err := gui.k8sClient.ListNamespace()
	if err != nil {
		return err
	}
	gui.data.NamespaceData = ns
	return nil
}

func (gui *Gui) reRenderNamespace() error {
	nsView := gui.getNamespaceView()
	if nsView == nil {
		return nil
	}

	if len(gui.data.NamespaceData) == 0 {
		return nil
	}

	gui.data.nsMux.RLock()
	defer gui.data.nsMux.RUnlock()
	ns := gui.data.NamespaceData

	gui.g.Update(func(*gocui.Gui) error {
		nsView.Clear()
		for _, n := range ns {
			fmt.Fprintln(nsView, n.Name)
		}
		return nil
	})

	return nil
}

func (gui *Gui) handleNSKeyUp(g *gocui.Gui, v *gocui.View) error {
	gui.changeSelectedLine(&gui.panelStates.Namespace.SelectedLine, len(gui.data.NamespaceData), false)
	return gui.handleNSSelect(v)
}

func (gui *Gui) handleNSKeyDown(g *gocui.Gui, v *gocui.View) error {
	gui.changeSelectedLine(&gui.panelStates.Namespace.SelectedLine, len(gui.data.NamespaceData), true)
	return gui.handleNSSelect(v)
}

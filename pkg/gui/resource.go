package gui

import (
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) getResourceView() *gocui.View {
	v, _ := gui.g.View("resource")
	return v
}

func (gui *Gui) onResourceClick(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	// Do something
	return nil
}

func (gui *Gui) onResourceTabClick(tabIndex int) error {

	resourceView := gui.getResourceView()
	resourceView.TabIndex = tabIndex

	infoView := gui.getInfoView()
	switch getResourceTabs()[tabIndex] {
	case "pod":
		infoView.Tabs = getPodInfoTabs()
	case "job":
		infoView.Tabs = getJobInfoTabs()
	case "deploy":
		infoView.Tabs = getDeployInfoTabs()
	case "service":
		infoView.Tabs = getServiceInfoTabs()
	case "secret":
		infoView.Tabs = getSecretInfoTabs()
	case "configMap":
		infoView.Tabs = getConfigMapInfoTabs()
	}

	return nil

}

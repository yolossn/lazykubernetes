package gui

import (
	"github.com/yolossn/lazykubernetes/pkg/utils"

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

func (gui *Gui) reRenderResource() error {
	resourceView := gui.getResourceView()

	if resourceView == nil {
		return nil
	}

	resources, err := gui.k8sClient.ListPods("default")

	if err != nil {
		return err
	}

	gui.g.Update(func(*gocui.Gui) error {
		resourceView.Clear()

				// make data for namespace tablewriter
		data := make([][]string, cap(resources))

		for x := 0; x < cap(resources); x++ {
			data[x] = make([]string, 3)
		}

		for i, n := range resources {
			data[i][0] = n.Name
		}

		utils.RenderTable(resourceView, data)

		return nil
	})

	return nil
}

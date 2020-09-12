package gui

import (
	"context"
	"fmt"
	"io"

	"github.com/jesseduffield/gocui"
	"gopkg.in/yaml.v2"
)

func (gui *Gui) getResourceView() *gocui.View {
	v, _ := gui.g.View("resource")
	return v
}

func (gui *Gui) onResourceClick(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	infoView := gui.getInfoView()
	// render pod
	switch getResourceTabs()[gui.panelStates.Resource.TabIndex] {
	case "pod":
		infoView.Tabs = getPodInfoTabs()
		return gui.handlePodSelect(v)
	}
	// podSelected := gui.FindSelectedLine(v, len(gui.data.PodData))
	// infoView.Clear()
	// pod := gui.data.PodData[podSelected]

	// Description 😘
	// data, err := gui.k8sClient.DescribePod(pod.Namespace, pod.Name)
	// if err != nil {
	// 	return err
	// }

	// output, err := yaml.Marshal(data)
	// if err != nil {
	// 	return err
	// }

	// fmt.Fprintln(infoView, string(output))

	// gui.g.Update(func(*gocui.Gui) error {
	// 	req := gui.k8sClient.StreamPodLogs(pod.Namespace, pod.Name)
	// 	ctx := context.TODO()
	// 	readCloser, err := req.Stream(ctx)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	infoView.Clear()
	// 	infoView.Autoscroll = true
	// 	go func() {
	// 		for {
	// 			io.Copy(infoView, readCloser)
	// 		}
	// 	}()
	// 	return nil
	// })

	return nil
}

func (gui *Gui) handlePodSelect(v *gocui.View) error {

	// Find Selected Pod
	podSelected := gui.FindSelectedLine(v, len(gui.data.PodData))
	pod := gui.data.PodData[podSelected]

	infoView := gui.getInfoView()

	// Find the tab in info panel
	switch getPodInfoTabs()[gui.panelStates.Resource.TabIndex] {
	case "logs":
		infoView.Clear()
		gui.g.Update(func(*gocui.Gui) error {
			ctx := context.TODO()
			req := gui.k8sClient.StreamPodLogs(pod.Namespace, pod.Name)
			readCloser, err := req.Stream(ctx)
			if err != nil {
				fmt.Println(err)
			}
			infoView.Clear()
			infoView.Autoscroll = true
			go func() {
				for {
					io.Copy(infoView, readCloser)
				}
			}()
			return nil
		})
	case "description":
		infoView.Clear()
		data, err := gui.k8sClient.DescribePod(pod.Namespace, pod.Name)
		if err != nil {
			return err
		}

		output, err := yaml.Marshal(data)
		if err != nil {
			return err
		}

		fmt.Fprintln(infoView, string(output))
	}
	return nil
}

func (gui *Gui) onResourceTabClick(tabIndex int) error {

	resourceView := gui.getResourceView()
	resourceView.TabIndex = tabIndex

	gui.panelStates.Resource.TabIndex = tabIndex
	infoView := gui.getInfoView()
	switch gui.getCurrentResourceTab() {
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
	rsView := gui.getResourceView()
	if rsView == nil {
		return nil
	}

	ns := gui.getCurrentNS()
	switch getResourceTabs()[gui.panelStates.Resource.TabIndex] {
	case "pod":
		gui.setPods(ns)
		return gui.renderPods()
	}

	return nil
}

func (gui *Gui) getCurrentResourceTab() string {
	return getResourceTabs()[gui.panelStates.Resource.TabIndex]
}

func (gui *Gui) setPods(namespace string) {
	gui.data.rsMux.Lock()
	defer gui.data.rsMux.Unlock()

	pods, err := gui.k8sClient.ListPods(namespace)
	if err != nil {

	}
	gui.data.PodData = pods
}

func (gui *Gui) renderPods() error {

	rsView := gui.getResourceView()
	if rsView == nil {
		return nil
	}

	gui.data.rsMux.RLock()
	defer gui.data.rsMux.RUnlock()

	rsView.Clear()
	pods := gui.data.PodData
	for _, pod := range pods {
		fmt.Fprintln(rsView, pod.Name)
	}
	return nil
}

func (gui *Gui) WatchPods() error {
	_ = gui.reRenderResource()
	// TODO: Handle error
	event, _ := gui.k8sClient.WatchPods("")
	for {
		_ = <-event.ResultChan()
		if gui.getCurrentResourceTab() != "pod" {
			continue
		}
		_ = gui.reRenderNamespace()
	}
}

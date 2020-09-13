package gui

import (
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/yolossn/lazykubernetes/pkg/utils"
	duration "k8s.io/apimachinery/pkg/util/duration"
)

func (gui *Gui) getNodeInfoView() *gocui.View {
	v, _ := gui.g.View("node")
	return v
}

func (gui *Gui) reRenderNodeInfo() error {

	nodeView := gui.getNodeInfoView()
	if nodeView == nil {
		return nil
	}

	if gui.getNSCount() == 0 {
		return nil
	}

	nodes, err := gui.k8sClient.ListNode()
	if err != nil {
		return err
	}

	nodeView.Clear()

	data := make([][]string, cap(nodes))

	for x := 0; x < cap(nodes); x++ {
		data[x] = make([]string, 4)
	}

	headers := []string{"NAME", "STATUS", "VERSION", "AGE"}
	for i, n := range nodes {
		data[i][0] = n.Name
		data[i][1] = n.Status
		data[i][2] = n.Version
		data[i][3] = duration.HumanDuration(time.Since(n.CreatedAt))
	}

	utils.RenderTable(nodeView, data, headers)

	return nil
}

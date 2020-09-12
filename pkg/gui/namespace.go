package gui

import (
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/olekukonko/tablewriter"

	duration "k8s.io/apimachinery/pkg/util/duration"
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

	return nil
}

func (gui *Gui) reRenderNamespace() error {
	nsView := gui.getNamespaceView()
	if nsView == nil {
		return nil
	}

	ns, err := gui.k8sClient.ListNamespace()
	if err != nil {
		return err
	}

	gui.g.Update(func(*gocui.Gui) error {
		nsView.Clear()

		// make data for namespace tablewriter
		data := make([][]string, cap(ns))

		for x := 0; x < cap(ns); x++ {
			data[x] = make([]string, 3)
		}

		for i, n := range ns {
			data[i][0] = n.Name
			data[i][1] = n.Status
			data[i][2] = duration.HumanDuration(time.Since(n.CreatedAt))
		}

		// Use table writer to render the data into view
		// https://github.com/olekukonko/tablewriter#example-10---set-nowhitespace-and-tablepadding-option
		table := tablewriter.NewWriter(nsView)
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t")
		table.SetNoWhiteSpace(true)
		table.AppendBulk(data)
		table.Render()

		return nil
	})

	return nil
}

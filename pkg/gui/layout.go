package gui

import (
	"github.com/jesseduffield/gocui"

	"github.com/yolossn/lazykubernetes/pkg/constants"
)

func (gui *Gui) layout(g *gocui.Gui) error {

	termWidth, termHeight := g.Size()

	// minimum size
	minimumHeight := 9
	minimumWidth := 10

	if termHeight < minimumHeight || termWidth < minimumWidth {
		v, err := g.SetView("limit", 0, 0, termWidth-1, termHeight-1, 0)
		if err != nil {
			if err.Error() != "unknown view" {
				return err
			}
			v.Title = constants.GetConstants().NotEnoughSpace
			v.Wrap = true
			_, _ = g.SetViewOnTop("limit")
		}
		return nil
	}

	_, _ = g.SetViewOnBottom("limit")
	g.DeleteView("limit")

	unitHeight := termHeight / 10

	leftColumnWidth := termWidth / 4

	if clusterInfoView, err := g.SetView("cluster-info", 0, 0, leftColumnWidth, unitHeight, gocui.BOTTOM|gocui.RIGHT); err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		clusterInfoView.Title = "cluster-info"
		clusterInfoView.Highlight = true
	}

	// namespaceViewHeight := termHeight - unitHeight - 1
	namespaceViewHeight := unitHeight * 2
	namespaceView, err := g.SetViewBeneath("namespace", "cluster-info", namespaceViewHeight)
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		namespaceView.Title = "namespace"
		namespaceView.Highlight = true
	}

	resourceViewHeight := unitHeight * 2
	if resourceView, err := g.SetView("resource", leftColumnWidth+1, 0, termWidth-1, resourceViewHeight, 0); err != nil {
		if err.Error() != "unknown view" {
			return err
		}
		resourceView.Tabs = getResourceTabs()
		resourceView.Highlight = true
	}

	infoViewHeight := termHeight - unitHeight*2 - 1
	_, err = g.SetViewBeneath("info", "resource", infoViewHeight)
	if err != nil {
		if err.Error() != "unknown view" {
			return err
		}
	}

	return nil
}

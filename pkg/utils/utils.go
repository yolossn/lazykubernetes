package utils

import (
	"io"

	"github.com/common-nighthawk/go-figure"
	"github.com/olekukonko/tablewriter"
)

func RenderTable(view io.Writer, data [][]string, header []string) {
	// Use table writer to render the data into view
	// https://github.com/olekukonko/tablewriter#example-10---set-nowhitespace-and-tablepadding-option

	table := tablewriter.NewWriter(view)
	table.SetHeader(header)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("   ")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data)
	table.Render()
}

func GetLazykubernetesArt() string {
	myFigure := figure.NewFigure("lazykubernetes", "stop", true)
	return myFigure.String()
}

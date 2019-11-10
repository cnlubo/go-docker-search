package main

import (
	"github.com/olekukonko/tablewriter"
)

// Display use to output something on screen with table format.
type DisplayTable struct {
	w *tablewriter.Table
}

func (d *DisplayTable) Render() {
	d.w.Render()
}

func (d *DisplayTable) SetHeader(head []string) {
	d.w.SetHeader(head)
}
func (d *DisplayTable) AppendBulk(rows [][]string) {
	d.w.AppendBulk(rows)
}

func (d *DisplayTable) SetHeaderColor(colors ...tablewriter.Colors) {

	var headerColors []tablewriter.Colors
	for _, color := range colors {
		headerColors = append(headerColors, tablewriter.Colors(color))
	}
	d.w.SetHeaderColor(headerColors...)
}

func (d *DisplayTable) SetColumnColor(colors ...tablewriter.Colors) {
	var columnColors []tablewriter.Colors
	for _, color := range colors {
		columnColors = append(columnColors, tablewriter.Colors(color))
	}
	d.w.SetColumnColor(columnColors...)
}
func (d *DisplayTable) SetBorders(border tablewriter.Border) {
	d.w.SetBorders(border)
}

func (d *DisplayTable) SetBorder(border bool) {
	d.w.SetBorder(border)
}
func (d *DisplayTable) SetRowLine(line bool) {
	d.w.SetRowLine(line)
}
func (d *DisplayTable) SetColumnSeparator(sep string) {
	d.w.SetColumnSeparator(sep)
}
func (d *DisplayTable) SetCenterSeparator(sep string) {
	d.w.SetCenterSeparator(sep)
}
func (d *DisplayTable) SetColumnAlignment(keys []int) {
	d.w.SetColumnAlignment(keys)
}

package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// ResponsiveTable wraps a table and automatically resizes columns proportionally
type ResponsiveTable struct {
	widget.BaseWidget
	table             *widget.Table
	columnProportions []float32 // Proportions for each column (should sum to 1.0)
	minWidth          float32   // Minimum total width before applying proportions
	padding           float32   // Reserved space for padding and separators
}

// NewResponsiveTable creates a responsive table with proportional column widths
// columnProportions: array of proportions for each column (e.g., [0.25, 0.5, 0.25] for 25%/50%/25%)
// minWidth: minimum total width before scaling kicks in (default 300)
// padding: space to reserve for padding/separators (default 20)
func NewResponsiveTable(table *widget.Table, columnProportions []float32, minWidth, padding float32) *ResponsiveTable {
	if minWidth <= 0 {
		minWidth = 300
	}
	if padding < 0 {
		padding = 20
	}

	r := &ResponsiveTable{
		table:             table,
		columnProportions: columnProportions,
		minWidth:          minWidth,
		padding:           padding,
	}
	r.ExtendBaseWidget(r)
	return r
}

func (r *ResponsiveTable) CreateRenderer() fyne.WidgetRenderer {
	return &responsiveTableRenderer{
		responsive: r,
		table:      r.table,
	}
}

type responsiveTableRenderer struct {
	responsive *ResponsiveTable
	table      *widget.Table
}

func (rr *responsiveTableRenderer) Layout(size fyne.Size) {
	rr.table.Resize(size)
	rr.responsive.updateColumnWidths(size.Width)
}

func (rr *responsiveTableRenderer) MinSize() fyne.Size {
	return rr.table.MinSize()
}

func (rr *responsiveTableRenderer) Refresh() {
	rr.table.Refresh()
}

func (rr *responsiveTableRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{rr.table}
}

func (rr *responsiveTableRenderer) Destroy() {}

// updateColumnWidths recalculates and applies proportional column widths
func (r *ResponsiveTable) updateColumnWidths(width float32) {
	if width < r.minWidth {
		return
	}

	availableWidth := width - r.padding
	if availableWidth < 0 {
		availableWidth = 0
	}
	// Don't enforce minWidth - allow columns to shrink to fit available space
	// This prevents horizontal scrolling when container is narrower than minWidth

	for i, proportion := range r.columnProportions {
		columnWidth := availableWidth * proportion
		r.table.SetColumnWidth(i, columnWidth)
	}
}

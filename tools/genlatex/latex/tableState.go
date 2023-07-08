package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/stylesheet"
	"math"
)

// tableState used to keep track of which columns are in use
type tableState struct {
	// Number of columns
	cols int
	// Column inactive state - set if a column is in a rowspan
	// 0 means it's active and expects a value.
	// Hidden columns get this set to MaxInt
	colRowSpan []int
	// Column definition cache
	colDefs []*stylesheet.ColumnSpec
}

func newTableState(cols int, table *stylesheet.Table, ctx context.Context) context.Context {
	ts := &tableState{
		cols:    cols,
		colDefs: table.ColumnSpec,
	}
	ts.reset()
	return context.WithValue(ctx, tableStateKey, ts)
}

func tableStateFromContext(ctx context.Context) *tableState {
	v := ctx.Value(tableStateKey)
	if v != nil {
		return v.(*tableState)
	}
	return nil
}

// reset() called at end of thead and tbody
func (ts *tableState) reset() {
	// init colRowSpan so hidden columns are ignored, all others are initially 0
	for col := 0; col < ts.cols; col++ {
		rowSpan := 0
		if ts.getCol(col).Hidden {
			rowSpan = math.MaxInt
		}
		ts.colRowSpan = append(ts.colRowSpan, rowSpan)
	}
}

func (ts *tableState) getCol(col int) *stylesheet.ColumnSpec {
	l := len(ts.colDefs)
	if col >= l {
		col = l - 1
	}
	return ts.colDefs[col]
}

// adjustCellNumber increments cell if it's on a column that's part of a row span
func (ts *tableState) adjustCellNumber(col int) int {
	for ; col < ts.cols; col++ {
		if ts.colRowSpan[col] > 0 {
			ts.colRowSpan[col]--
		}
	}
	return col
}

// incrementCellNumber adds colSpan to cell accounting for any hidden Columns
func (ts *tableState) incrementCellNumber(col, colSpan int) int {
	for ; colSpan > 0; colSpan-- {
		// If hidden then skip this column
		//for col < ts.cols && ts.getCol(col).Hidden {
		//	col++
		//}

		// Increment col
		col++
	}
	return col
}

// fixRowSpan ensures we have an entry for the specified column
func (ts *tableState) fixRowSpan(col int) {
	col++
	for len(ts.colRowSpan) < col {
		ts.colRowSpan = append(ts.colRowSpan, 0)
	}
}

// setRowSpan marks a column has a row span
func (ts *tableState) setRowSpan(col, rowSpan int) {
	ts.fixRowSpan(col)
	// -1 because we are already in this column
	ts.colRowSpan[col] = rowSpan - 1
}

func (ts *tableState) decRowSpan(col int) bool {
	ts.fixRowSpan(col)
	r := ts.colRowSpan[col] > 0
	if r {
		ts.colRowSpan[col] = ts.colRowSpan[col] - 1
	}
	return r
}

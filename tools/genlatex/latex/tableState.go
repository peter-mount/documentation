package latex

import (
	"context"
	"github.com/peter-mount/documentation/tools/genlatex/stylesheet"
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
	// Table model
	table *Table
}

func newTableState(tType string, cols int, table *stylesheet.Table, ctx context.Context) context.Context {
	ts := &tableState{
		cols:    cols,
		colDefs: table.ColumnSpec,
		table: &Table{
			Type:  tType,
			Table: table,
		},
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

func (ts *tableState) reset() {
	ts.fixRowSpan(ts.cols - 1)
	for col := 0; col < ts.cols; col++ {
		ts.colRowSpan[col] = 0
	}
}

func (ts *tableState) getCol(col int) *stylesheet.ColumnSpec {
	l := len(ts.colDefs)
	if col >= l {
		col = l - 1
	}
	return ts.colDefs[col]
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

// decRowSpan will, if a cell is part of a rowspan, decrement the columns rowSpan counter by 1.
// The return bool will be true if the column is still in this state after the decrement.
func (ts *tableState) decRowSpan(col int) bool {
	if col < len(ts.colRowSpan) && ts.colRowSpan[col] > 0 {
		ts.colRowSpan[col] = ts.colRowSpan[col] - 1
		return true
	}
	return false
}

func (ts *tableState) finishRow(col int) {
	for col < len(ts.colRowSpan) {
		ts.decRowSpan(col)
		col++
	}
}

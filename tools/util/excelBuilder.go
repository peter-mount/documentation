package util

import (
  "context"
  "github.com/xuri/excelize/v2"
  "io"
)

type ExcelProvider interface {
  GetExcel() ExcelBuilder
  SetExcel(ExcelBuilder)
}

type ExcelBuilder func(context.Context, *excelize.File) error

func NewExcelBuilder() ExcelBuilder {
  return nil
}

func (a ExcelBuilder) Then(b ExcelBuilder) ExcelBuilder {
  if a == nil {
    return b
  }
  return func(ctx context.Context, file *excelize.File) error {
    if err := a(ctx, file); err != nil {
      return err
    }
    return b(ctx, file)
  }
}

func (a ExcelBuilder) After(b ExcelBuilder) ExcelBuilder {
  return b.Then(a)
}

func (a ExcelBuilder) FileHandler() FileHandler {
  return func(w io.Writer) error {
    f := excelize.NewFile()

    err := a(context.Background(), f)
    if err != nil {
      return err
    }

    f.DeleteSheet("Sheet1")

    return f.Write(w)
  }
}

// CellName returns the Excel cell name.
// Note: col & row start from 1
func CellName(col, row int, abs ...bool) string {
  r, _ := excelize.CoordinatesToCellName(col, row, abs...)
  return r
}

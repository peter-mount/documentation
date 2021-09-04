package util

import (
  "github.com/xuri/excelize/v2"
  "io"
)

type ExcelProvider interface {
  GetExcel() ExcelBuilder
  SetExcel(ExcelBuilder)
}

type ExcelBuilder func(*excelize.File) error

func NewExcelBuilder() ExcelBuilder {
  return func(_ *excelize.File) error {
    return nil
  }
}

func (a ExcelBuilder) Then(b ExcelBuilder) ExcelBuilder {
  return func(file *excelize.File) error {
    err := a(file)
    if err != nil {
      return err
    }
    return b(file)
  }
}

func (a ExcelBuilder) After(b ExcelBuilder) ExcelBuilder {
  //return b.Then(a)
  return func(file *excelize.File) error {
    err := b(file)
    if err != nil {
      return err
    }
    return a(file)
  }
}

func (a ExcelBuilder) FileHandler() FileHandler {
  return func(w io.Writer) error {
    f := excelize.NewFile()

    err := a(f)
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

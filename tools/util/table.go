package util

import (
  "github.com/xuri/excelize/v2"
  "time"
)

// Table contains the definition of a reference table which is then written to a csv or json file.
// Multiple Table's are also written to a single Excel xlsx file, one per Sheet.
type Table struct {
  Title     string         // Title of the report
  Columns   []string       // Column definitions
  RowCount  int            // Number of data rows
  GetRow    TableGetRow    // Get record for a specific row
  Transform TableTransform // Transfer a record to a row of strings
}

type TableGetRow func(r int) interface{}
type TableTransform func(interface{}) []interface{}

func (t Table) ForEachRow(h func(int, int, []interface{}) error) error {
  for r := 0; r < t.RowCount; r++ {
    err := h(r, t.RowCount, t.Transform(t.GetRow(r)))
    if err != nil {
      return err
    }
  }
  return nil
}

func WithTable() TableHandler {
  return func(table *Table) error {
    return nil
  }
}

type TableHandler func(*Table) error

func (a TableHandler) Then(b TableHandler) TableHandler {
  return func(table *Table) error {
    err := a(table)
    if err != nil {
      return err
    }
    return b(table)
  }
}

func (a TableHandler) ForEach(handlers ...TableHandler) TableHandler {
  return func(table *Table) error {
    if a != nil {
      err := a(table)
      if err != nil {
        return err
      }
    }

    for _, h := range handlers {
      err := h(table)
      if err != nil {
        return err
      }
    }

    return nil
  }
}

func (a TableHandler) Do(t *Table) error {
  return a(t)
}

func (a TableHandler) AsCSV(fileName string, fileTime time.Time) TableHandler {
  return a.Then(func(t *Table) error {
    return NewCSVBuilder().
      Headings(t.Columns...).
      ImportFrom(t).
      FileHandler().
      Write(fileName, fileTime)
  })
}

func (a TableHandler) AsExcel(p ExcelProvider) TableHandler {
  return a.Then(func(t *Table) error {
    p.SetExcel(p.GetExcel().Then(func(f *excelize.File) error {
      f.NewSheet(t.Title)
      for i, c := range t.Columns {
        _ = f.SetCellValue(t.Title, CellName(i+1, 1), c)
      }

      return t.ForEachRow(func(r int, rc int, v []interface{}) error {
        for i, c := range v {
          _ = f.SetCellValue(t.Title, CellName(i+1, r+2), c)
        }
        return nil
      })
    }))
    return nil
  })
}

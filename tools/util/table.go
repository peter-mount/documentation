package util

import (
  "context"
  "github.com/xuri/excelize/v2"
  "time"
)

// Table contains the definition of a reference table which is then written to a csv or json file.
// Multiple Table's are also written to a single Excel xlsx file, one per Sheet.
type Table struct {
  Title          string         // Title of the report
  Columns        []string       // Column definitions
  RowCount       int            // Number of data rows
  GetRow         TableGetRow    // Get record for a specific row
  Transform      TableTransform // Transfer a record to a row of strings
  widths         map[int]int    // Map of column widths
  dataStyleId    int            // Excel styleId for data
  headingStyleId int            // Excel styleId for headings
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

// Do passes a Table to a TableHandler
func (a TableHandler) Do(t *Table) error {
  if a != nil {
    return a(t)
  }
  return nil
}

// AsCSV will create a CSV file based on this Table
func (a TableHandler) AsCSV(fileName string, fileTime time.Time) TableHandler {
  return a.Then(func(t *Table) error {
    return NewCSVBuilder().
      Headings(t.Columns...).
      ImportFrom(t).
      FileHandler().
      Write(fileName, fileTime)
  })
}

func (t *Table) setCell(f *excelize.File, c, r int, v interface{}) {
  axis := CellName(c, r)

  _ = f.SetCellValue(t.Title, axis, v)

  vs, err := f.GetCellValue(t.Title, axis)
  if err == nil {
    l := len(vs)
    if el, exists := t.widths[c]; !exists || el < l {
      t.widths[c] = l
    }
  }
}

// AsExcel will create a new Sheet in an Excel Workbook based on this Table
func (a TableHandler) AsExcel(p ExcelProvider) TableHandler {
  return a.Then(func(t *Table) error {
    return p.BuildExcel(func(builder ExcelBuilder) ExcelBuilder {
      return builder.Then(func(ctx context.Context, f *excelize.File) error {
        t.widths = make(map[int]int)

        f.NewSheet(t.Title)

        for i, c := range t.Columns {
          t.setCell(f, i+1, 1, c)
        }

        _ = t.ForEachRow(func(r int, rc int, v []interface{}) error {
          for i, c := range v {
            t.setCell(f, i+1, r+2, c)
          }
          return nil
        })

        minC, maxC := -1, -1
        for c, w := range t.widths {
          cs, _ := excelize.ColumnNumberToName(c)
          _ = f.SetColWidth(t.Title, cs, cs, float64(w))

          if minC < 1 || minC > c {
            minC = c
          }
          if maxC < c {
            maxC = c
          }
        }

        SetCellStyle(f, &t.headingStyleId, t.Title, minC, 1, maxC, 1, true, true)
        SetCellStyle(f, &t.dataStyleId, t.Title, minC, 2, maxC, 2+t.RowCount, false, false)

        return nil
      })
    })
  })
}

func SetCellStyle(f *excelize.File, id *int, sheet string, c1, r1, c2, r2 int, bold, header bool) {
  if *id == 0 {
    s := `{"font":{`

    if bold {
      s = s + `"bold":true,`
    }

    s = s + `"color":"000000",`
    s = s + `"family":"Courier","size":10}`

    if header {
      s = s + `,"border":[{"type":"bottom","color":"000000","style":1}]`
      s = s + `,"fill":{"type":"pattern","color":["aaaaaa"],"pattern":1}`
    }

    s = s + `}`

    *id, _ = f.NewStyle(s)
  }

  axis1, _ := excelize.CoordinatesToCellName(c1, r1)
  axis2, _ := excelize.CoordinatesToCellName(c2, r2)
  _ = f.SetCellStyle(sheet, axis1, axis2, *id)
}

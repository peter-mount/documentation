package util

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

func (a TableHandler) Do(t *Table) error {
  /*  for i:=0;i<t.RowCount;i++ {
        r:=t.Transform(t.GetRow(i))
        for c
      }
  */return a(t)
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

func (t Table) AsCSV() CSVBuilder {
  return NewCSVBuilder().
    Headings(t.Columns...).
    ImportFrom(t)
}

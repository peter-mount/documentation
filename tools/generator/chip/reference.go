package chip

import (
  "fmt"
  "github.com/peter-mount/documentation/tools/hugo"
  "github.com/peter-mount/documentation/tools/util"
)

// Generates the chip reference tables.
// This is usually only used for the chipref book
func (c *Chip) chipReferenceTables(book *hugo.Book) error {
  c.generator.AddTask(func() error {
    return c.chipReferenceTablesTask(book)
  })
  return nil
}

func (c *Chip) chipReferenceTablesTask(book *hugo.Book) error {
  return c.chips.ForEachCategory(func(cat string) error {
    return util.WithTable().
      AsCSV(book.StaticPath(cat+".csv"), book.Modified()).
      AsExcel(c.excel.Get(book.ID, book.Modified())).
      Do(c.chipReferenceTable(book, cat))
  })
}

func (c *Chip) chipReferenceTable(book *hugo.Book, cat string) *util.Table {

  defs := c.chips.DefinitionNames(cat)

  // Get max pin count for this category
  maxPins := 0
  _ = defs.ForEach(func(name string) error {
    d := c.chips.Get(cat, name)
    if d != nil && d.PinCount > maxPins {
      maxPins = d.PinCount
    }
    return nil
  })

  t := &util.Table{
    Title: cat,
    Columns: []string{
      "Name",
      "Title",
      "Category",
      "SubCategory",
      "Label",
      "Sub Label",
      "Type",
      "Pins",
    },
    RowCount: len(defs),
    GetRow: func(r int) interface{} {
      return c.chips.Get(cat, defs[r])
    },
  }

  // Add pin columns up to MaxPins
  for pin := 1; pin <= maxPins; pin++ {
    t.Columns = append(t.Columns, fmt.Sprintf("Pin %d", pin))
  }

  t.Transform = func(i interface{}) []interface{} {
    d := i.(*Definition)

    var a []interface{}
    a = append(a,
      d.Name,
      d.Title,
      d.Category,
      d.SubCategory,
      d.Label,
      d.SubLabel,
      d.Type,
      d.PinCount,
    )

    // Add chip pin definitions
    for i := 1; i <= d.PinCount; i++ {
      a = append(a, d.Pins[i])
    }
    // Pad to maxPins
    if d.PinCount < maxPins {
      for i := d.PinCount + 1; i <= maxPins; i++ {
        a = append(a, "")
      }
    }

    return a
  }

  return t
}

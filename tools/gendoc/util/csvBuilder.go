package util

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type CSVBuilder func(*csv.Writer) error

func NewCSVBuilder() CSVBuilder {
	return func(w *csv.Writer) error {
		return nil
	}
}

func (a CSVBuilder) Then(b CSVBuilder) CSVBuilder {
	return func(w *csv.Writer) error {
		err := a(w)
		if err != nil {
			return err
		}
		return b(w)
	}
}

func (a CSVBuilder) Headings(h ...string) CSVBuilder {
	return func(w *csv.Writer) error {
		err := w.Write(h)
		if err != nil {
			return err
		}
		return a(w)
	}
}

func (a CSVBuilder) DOS() CSVBuilder {
	return func(w *csv.Writer) error {
		w.UseCRLF = true
		return a(w)
	}
}

func (a CSVBuilder) Unix() CSVBuilder {
	return func(w *csv.Writer) error {
		w.UseCRLF = false
		return a(w)
	}
}

func (a CSVBuilder) Separator(comma rune) CSVBuilder {
	return func(w *csv.Writer) error {
		w.Comma = comma
		return a(w)
	}
}

type CSVSource interface {
	ForEachRow(h func(int, int, []interface{}) error) error
}

func InterfaceToString(row []interface{}) []string {
	var a []string

	for _, e := range row {
		switch v := e.(type) {
		case string:
			a = append(a, v)
		case int:
			a = append(a, strconv.Itoa(v))
		case bool:
			a = append(a, strconv.FormatBool(v))
		default:
			a = append(a, fmt.Sprintf("%v", v))
		}
	}

	return a
}

func (a CSVBuilder) ImportFrom(src CSVSource) CSVBuilder {
	return func(w *csv.Writer) error {
		err := a(w)
		if err != nil {
			return err
		}

		return src.ForEachRow(func(_ int, _ int, row []interface{}) error {
			return w.Write(InterfaceToString(row))
		})
	}
}

func (a CSVBuilder) FileHandler() FileHandler {
	return func(w io.Writer) error {
		nw := csv.NewWriter(w)

		err := a(nw)
		if err != nil {
			return err
		}

		nw.Flush()
		return nw.Error()
	}
}

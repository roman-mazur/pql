package visualize

import (
	"fmt"
	"io"

	"rmazur.io/pql/data"
)

func Table(ds data.Set, limit int, out io.Writer) error {
	if rs, reusable := ds.(data.ReusableSet); reusable {
		_, _ = fmt.Fprintln(out, "COUNT:", rs.Count())
	}

	cols, err := ds.Columns()
	if err != nil {
		return err
	}

	w := len(cols)
	holder := make([]string, w)
	row := make([]interface{}, w)
	for i := range row {
		holder[i] = cols[i]
		row[i] = &holder[i]
	}
	printRow(out, row) // Print the header.

	for i := 0; i < limit && ds.Next(); i++ {
		if err := ds.Scan(row...); err != nil {
			return err
		}
		printRow(out, row)
	}
	return nil
}

func printRow(out io.Writer, values []interface{}) {
	for i, v := range values {
		if s, ok := v.(*string); ok {
			_, _ = fmt.Fprint(out, *s)
		} else {
			_, _ = fmt.Fprintf(out, "%v", v)
		}
		if i != len(values)-1 {
			_, _ = fmt.Fprint(out, "\t")
		}
	}
	_, _ = fmt.Fprintln(out)
}

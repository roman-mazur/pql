package visualize

import (
	"fmt"

	"rmazur.io/pql/data"
)

func Table(ds data.Set, limit int) error {
	if rs, reusable := ds.(data.ReusableSet); reusable {
		fmt.Println("COUNT:", rs.Count())
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
	printRow(row) // Print the header.

	for i := 0; i < limit && ds.Next(); i++ {
		if err := ds.Scan(row...); err != nil {
			return err
		}
		printRow(row)
	}
	return nil
}

func printRow(values []interface{}) {
	for i, v := range values {
		if s, ok := v.(*string); ok {
			fmt.Print(*s)
		} else {
			fmt.Printf("%v", v)
		}
		if i != len(values)-1 {
			fmt.Print("\t")
		}
	}
	fmt.Println()
}

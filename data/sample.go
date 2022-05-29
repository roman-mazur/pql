package data

import (
	"fmt"
	"reflect"
)

type value struct {
	V float64
	L string
}

type ChartedValues struct {
	v       []value
	columns []string

	vi, li int
	c      int
}

func ReadAll(d Set, vi, li int) (*ChartedValues, error) {
	col, _ := d.Columns()
	cv := &ChartedValues{
		columns: col,
		vi:      vi,
		li:      li,
	}
	cv.Reset()

	var (
		v float64
		l string
	)
	args := ScanArgs(&v, &l, vi, li, len(col))

	for d.Next() {
		if err := d.Scan(args...); err != nil {
			return nil, err
		}
		cv.v = append(cv.v, value{V: v, L: l})
	}
	return cv, nil
}

func (v *ChartedValues) Reset() {
	v.c = -1
}

func (v *ChartedValues) Next() bool {
	v.c++
	return v.c < len(v.v)
}

func (v *ChartedValues) Err() error {
	return nil
}

func (v *ChartedValues) Scan(dest ...interface{}) error {
	vt := dest[v.vi]
	if dv, ok := vt.(*float64); ok {
		*dv = v.v[v.c].V
	} else {
		return fmt.Errorf("unexpected dest for the value (index %d, %s)", v.vi, reflect.TypeOf(vt))
	}
	lt := dest[v.li]
	if dl, ok := lt.(*string); ok {
		*dl = v.v[v.c].L
	}
	return nil
}

func (v *ChartedValues) Columns() ([]string, error) {
	return v.columns, nil
}

func (v *ChartedValues) Close() error {
	return nil
}

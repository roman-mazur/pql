package data

import (
	"fmt"
	"reflect"
	"strconv"
)

type value struct {
	V float64
	L string
}

type ChartedValues struct {
	arrayValues

	v      []value
	vi, li int
}

func ReadAll(d Set, vi, li int) (*ChartedValues, error) {
	col, _ := d.Columns()
	cv := &ChartedValues{
		arrayValues: arrayValues{
			columns: col,
		},
		vi: vi,
		li: li,
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
	cv.l = len(cv.v)
	return cv, nil
}

func (v *ChartedValues) Scan(dest ...interface{}) error {
	vt := dest[v.vi]
	switch dv := vt.(type) {
	case *float64:
		*dv = v.v[v.c].V
	case *string:
		*dv = strconv.FormatFloat(v.v[v.c].V, 'f', -1, 64)
	default:
		return fmt.Errorf("unexpected dest for the value (index %d, %s)", v.vi, reflect.TypeOf(vt))
	}

	lt := dest[v.li]
	switch dl := lt.(type) {
	case *float64:
		lv, err := strconv.ParseFloat(v.v[v.c].L, 64)
		if err != nil {
			return fmt.Errorf("can't parse numerical label (index %d, %s)", v.vi)
		}
		*dl = lv
	case *string:
		*dl = v.v[v.c].L
	default:
		return fmt.Errorf("unexpected dest for the label (index %d, %s)", v.vi, reflect.TypeOf(vt))
	}

	return nil
}

type arrayValues struct {
	columns []string
	c       int
	l       int
}

func (av *arrayValues) reusable() {}

func (av *arrayValues) Columns() ([]string, error) {
	return av.columns, nil
}

func (av *arrayValues) Count() int {
	return av.l
}

func (av *arrayValues) Close() error {
	return nil
}

func (av *arrayValues) Reset() {
	av.c = -1
}

func (av *arrayValues) Next() bool {
	av.c++
	return av.c < av.l
}

func (av *arrayValues) Err() error {
	return nil
}

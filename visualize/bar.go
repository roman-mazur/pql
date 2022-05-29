package visualize

import (
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"rmazur.io/pql/data"
)

func readValues(s data.Set, vi, li, size int) (plotter.Valuer, []string, error) {
	const n = 20
	values := make(plotter.Values, 0, n)
	labels := make([]string, 0, n)

	var (
		l string
		v float64
	)
	args := data.ScanArgs(&v, &l, vi, li, size)

	for s.Next() {
		err := s.Scan(args...)
		if err != nil {
			return nil, nil, err
		}
		values = append(values, v)
		if size > 1 {
			labels = append(labels, l)
		} else {
			labels = append(labels, strconv.Itoa(len(values)))
		}
	}

	if err := s.Err(); err != nil {
		return nil, nil, err
	}
	return values, labels, nil
}

func Bar(name string, s data.Set, x, y string) error {
	c, err := s.Columns()
	if err != nil {
		return err
	}
	vi, li := data.Indexes(c, x, y)

	v, l, err := readValues(s, vi, li, len(c))
	if err != nil {
		return err
	}
	bc, err := plotter.NewBarChart(v, windowWidth)
	if err != nil {
		return err
	}

	p := plot.New()
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(bc, plotter.NewGrid())
	p.NominalX(l...)

	display(name, p)
	return nil
}

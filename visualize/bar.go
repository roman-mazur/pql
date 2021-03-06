package visualize

import (
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
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

	width := windowWidth / font.Length(len(l)+1)
	if width > 100 {
		width = 100
	}
	bc, err := plotter.NewBarChart(v, width)
	if err != nil {
		return err
	}
	p := plot.New()

	lx, ly := x, y
	if lx == "" {
		lx = "X"
	}
	if ly == "" {
		ly = "Y"
	}
	p.X.Label.Text = lx
	p.Y.Label.Text = ly
	const labelSize = windowWidth / 30
	p.X.Label.TextStyle.Font.Size = labelSize
	p.Y.Label.TextStyle.Font.Size = labelSize

	p.Add(bc, plotter.NewGrid())
	p.NominalX(l...)
	p.X.Tick.Label.Font.Size = p.X.Label.TextStyle.Font.Size * 0.66
	p.Y.Tick.Label.Font.Size = p.Y.Label.TextStyle.Font.Size * 0.66

	display(name, p)
	return nil
}

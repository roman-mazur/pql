package main

import (
	"flag"
	"fmt"
	"image/color"
	"math"
	"os"
	"os/signal"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vggio"
)

func repl(init Command) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	repl := NewRepl("select * from testdata")

	go func() {
		<-interrupt
		fmt.Println("Bye!")
		for repl.CurrentCtx() != nil { // TODO: It's a race with the main REPL loop.
			repl.FinishCtx()
		}
		os.Exit(0)
	}()

	for cmd := init; cmd != nil; cmd = repl.next() {
		cmd.Perform(repl)
	}
	os.Exit(0)
}

type cmdList []Command

func (cl cmdList) Perform(repl *Repl) {
	for i := range cl {
		cl[i].Perform(repl)
	}
}

func main() {
	postgresIn := flag.String("postgres", "", "connection string for a Postgres db")
	osqueryIn := flag.String("osquery", "", "osqueryd extensions socket path")
	flag.Parse()

	init := cmdList{switchCmd("$")}
	if *postgresIn != "" {
		init = append(init, &connectCmd{Driver: "postgres", ConnStr: *postgresIn})
	} else if *osqueryIn != "" {
		path := *osqueryIn
		if path == "-" {
			path = "/var/osquery/osquery.em"
		}
		init = append(init, &connectCmd{Driver: "osquery", ConnStr: path})
	}

	go repl(init)

	win := app.NewWindow(app.Title("test"), app.Size(unit.Px(400), unit.Px(400)))
	defer win.Close()
	go loop(win)

	app.Main()
}

func loop(w *app.Window) error {
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			p := plot.New()
			p.Title.Text = "My title"
			p.X.Label.Text = "X"
			p.Y.Label.Text = "Y"

			quad := plotter.NewFunction(func(x float64) float64 {
				return x * x
			})
			quad.Color = color.RGBA{B: 255, A: 255}

			exp := plotter.NewFunction(func(x float64) float64 {
				return math.Pow(2, x)
			})
			exp.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
			exp.Width = vg.Points(2)
			exp.Color = color.RGBA{G: 255, A: 255}

			sin := plotter.NewFunction(func(x float64) float64 {
				return 10*math.Sin(x) + 50
			})
			sin.Dashes = []vg.Length{vg.Points(4), vg.Points(5)}
			sin.Width = vg.Points(4)
			sin.Color = color.RGBA{R: 255, A: 255}

			p.Add(quad, exp, sin)
			p.Legend.Add("x^2", quad)
			p.Legend.Add("2^x", exp)
			p.Legend.Add("10*sin(x)+50", sin)
			p.Legend.ThumbnailWidth = 0.5 * vg.Inch

			p.X.Min = 0
			p.X.Max = 10
			p.Y.Min = 0
			p.Y.Max = 100

			p.Add(plotter.NewGrid())

			gtx := layout.NewContext(&ops, e)
			const ww = 15 * vg.Centimeter
			cnv := vggio.New(gtx, ww, ww, vggio.UseDPI(96))
			p.Draw(draw.New(cnv))

			e.Frame(cnv.Paint())
		}
	}
}

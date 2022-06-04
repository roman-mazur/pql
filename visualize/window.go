package visualize

import (
	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vggio"
)

const (
	windowWidth = 30 * vg.Centimeter
	dpi         = 96
)

func display(name string, p *plot.Plot) {
	s := unit.Px(float32(windowWidth.Dots(dpi)))
	win := app.NewWindow(
		app.Title(name),
		app.Size(s, s),
	)

	go func() {
		defer win.Close()

		var ops op.Ops
	eloop:
		for e := range win.Events() {
			switch e := e.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				cnv := vggio.New(gtx, windowWidth, windowWidth, vggio.UseDPI(dpi))
				p.Draw(draw.New(cnv))
				e.Frame(&ops)

			case key.Event:
				switch e.Name {
				case key.NameEscape:
					break eloop
				}

			case system.DestroyEvent:
				break eloop
			}
		}
	}()
}

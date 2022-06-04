package main

import (
	"fmt"

	"rmazur.io/pql/data"
	"rmazur.io/pql/visualize"
)

type barCmd struct {
	X, Y string
}

func (bc *barCmd) Perform(repl *Repl) {
	if err := persistDs(repl, bc.X, bc.Y); err != nil {
		repl.MsgErr(err)
		return
	}

	cc := repl.CurrentCtx()
	err := visualize.Bar(cc.Name, cc.S, bc.X, bc.Y)
	if err != nil {
		repl.MsgErr(err)
	}
}

type tableCmd int

func (tc tableCmd) Perform(repl *Repl) {
	if err := persistDs(repl, "", ""); err != nil {
		repl.MsgErr(err)
		return
	}

	cc := repl.CurrentCtx()
	limit := int(tc)
	if limit == 0 {
		limit = 5
	}
	err := visualize.Table(cc.S, limit)
	if err != nil {
		repl.MsgErr(err)
	}
}

func persistDs(repl *Repl, x, y string) error {
	cc := repl.CurrentCtx()
	if cc.S == nil {
		return fmt.Errorf("no data in the current context")
	}

	cv, ok := cc.S.(data.ReusableSet)
	if !ok {
		col, err := cc.S.Columns()
		if err != nil {
			return err
		}
		vi, li := data.Indexes(col, x, y)
		cv, err = data.ReadAll(cc.S, vi, li)
		if err != nil {
			return err
		}
		_ = cc.S.Close()
	}
	cv.Reset()
	cc.S = cv
	return nil
}

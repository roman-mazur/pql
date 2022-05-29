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
	cc := repl.CurrentCtx()
	if cc.S == nil {
		repl.MsgErr(fmt.Errorf("no data in the current context"))
		return
	}

	if err := bc.persist(cc); err != nil {
		repl.MsgErr(err)
		return
	}

	err := visualize.Bar(cc.Name, cc.S, bc.X, bc.Y)
	if err != nil {
		repl.MsgErr(err)
	}
}

func (bc *barCmd) persist(cc *CmdCtx) error {
	cv, ok := cc.S.(*data.ChartedValues)
	if !ok {
		col, err := cc.S.Columns()
		if err != nil {
			return err
		}
		vi, li := data.Indexes(col, bc.X, bc.Y)
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

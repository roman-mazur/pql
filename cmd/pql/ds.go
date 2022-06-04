package main

import (
	"fmt"

	"rmazur.io/pql/data"
)

type connectCmd struct {
	Driver  string
	ConnStr string
}

func (c *connectCmd) Perform(repl *Repl) {
	cc := repl.CurrentCtx()
	if cc.Ds != nil {
		switchCmd("$" + cc.Name).Perform(repl)
	}

	var (
		ds  data.Source
		err error
	)
	if c.Driver == "osquery" {
		ds, err = data.NewOsQuerySource(c.ConnStr)
	} else {
		ds, err = data.OpenSQL(c.Driver, c.ConnStr)
	}

	if err != nil {
		repl.MsgErr(err)
		return
	}
	repl.CurrentCtx().Ds = ds
}

type queryCmd string

func (qc queryCmd) Perform(repl *Repl) {
	cc := repl.CurrentCtx()
	if cc.Ds == nil {
		repl.MsgErr(fmt.Errorf("connection is not established"))
		return
	}

	set, err := cc.Ds.Query(repl.CurrentCtx(), string(qc))
	if err != nil {
		repl.MsgErr(err)
		return
	}

	repl.CurrentCtx().UpdateSet(set)
}

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
	"rmazur.io/pql/data"
)

type CmdCtx struct {
	context.Context

	Name string
	Ds   data.Source

	S data.Set
}

func (ctx *CmdCtx) UpdateSet(s data.Set) {
	if ctx.S != nil {
		_ = ctx.S.Close()
	}
	ctx.S = s
}

type Command interface {
	Perform(ctx *Repl)
}

type CommandFunc func(repl *Repl)

func (f CommandFunc) Perform(repl *Repl) {
	f(repl)
}

type Repl struct {
	input  io.Reader
	output io.Writer

	bin *bufio.Reader
	t   *term.Terminal

	ctx map[string]*CmdCtx
	s   []string
	fs  []string // First commands to execute.
}

func NewRepl(cmds ...string) *Repl {
	return &Repl{
		ctx: make(map[string]*CmdCtx),
		fs:  cmds,
	}
}

func (r *Repl) out() io.Writer {
	if r.output == nil {
		if r.t != nil {
			return r.t
		}
		return os.Stdout
	}
	return r.output
}

func (r *Repl) in() *bufio.Reader {
	if r.input == nil {
		if r.t == nil {
			_, err := term.MakeRaw(int(os.Stdin.Fd()))
			if err == nil {
				r.t = term.NewTerminal(os.Stdin, "")
			} else {
				log.Println("Failed to switch to a raw terminal mode:", err)
				r.bin = bufio.NewReaderSize(os.Stdin, 2)
			}
		}
	} else if r.bin == nil {
		r.bin = bufio.NewReaderSize(r.input, 256)
	}
	return r.bin
}

func (r *Repl) push(name string) {
	r.s = append(r.s, name)
}

func (r *Repl) CurrentCtx() *CmdCtx {
	if len(r.s) == 0 {
		return nil
	}
	return r.ctx[r.s[len(r.s)-1]]
}

func (r *Repl) Switch(name string) {
	cc := r.CurrentCtx()
	r.push(name)
	if _, exists := r.ctx[name]; !exists {
		ctx := context.Background()
		if cc != nil {
			ctx = cc.Context
		}
		r.ctx[name] = &CmdCtx{
			Context: ctx,
			Name:    name,
		}
	}
}

func (r *Repl) MsgErr(err error) {
	_, _ = fmt.Fprintln(r.out(), "ERROR:", err)
}

func (r *Repl) FinishCtx() {
	if len(r.s) == 0 {
		panic("context stack is empty")
	}
	cc := r.CurrentCtx()

	if cc.Ds != nil {
		_ = cc.Ds.Close()
	}

	delete(r.ctx, cc.Name)
	r.s[len(r.s)-1] = ""
	r.s = r.s[:len(r.s)-1]
}

func (r *Repl) prompt() {
	cc := r.CurrentCtx()
	_, _ = fmt.Fprint(r.out(), cc.Name, "> ")
}

func (r *Repl) read() (string, error) {
	if len(r.fs) > 0 {
		res := r.fs[0]
		r.fs = r.fs[1:]
		return res, nil
	}

	if bin := r.in(); bin != nil {
		return bin.ReadString('\n')
	}

	// Working with a terminal.
	if r.t == nil {
		panic("invalid state: terminal is not initialized")
	}
	return r.t.ReadLine()
}

func (r *Repl) next() Command {
	if len(r.s) == 0 {
		return nil
	}

	r.prompt()

	line, err := r.read()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return errorCmd{err: err}
	}
	line = strings.TrimSpace(line)

	word := strings.SplitN(line, " ", 2)[0]
	cmd, err := makeCmd(word, line)
	if err != nil {
		cmd = errorCmd{err: err}
	}
	return cmd
}

func makeCmd(word, line string) (cmd Command, err error) {
	defer func() {
		if e := recover(); e != nil {
			cmd = nil
			err = fmt.Errorf("problems parsing command: %s", e)
		}
	}()

	switch strings.ToLower(word) {
	case "switch":
		args := strings.Fields(line)
		cmd = switchCmd(args[1])

	case "quit":
		cmd = CommandFunc(quit)

	case "connect":
		args := strings.SplitN(line, " ", 3)
		cmd = &connectCmd{Driver: args[1], ConnStr: args[2]}

	case "":
		cmd = CommandFunc(nop)

	case "bar":
		args := strings.Fields(line)
		x, y := "", ""
		if len(args) > 1 {
			x = args[1]
		}
		if len(args) > 2 {
			y = args[2]
		}
		cmd = &barCmd{X: x, Y: y}

	case "table":
		args := strings.Fields(line)
		limit := 0
		if len(args) > 1 {
			limit, _ = strconv.Atoi(args[1])
		}
		cmd = tableCmd(limit)

	default:
		cmd = queryCmd(line)
	}
	return
}

type switchCmd string

func (s switchCmd) Perform(repl *Repl) {
	repl.Switch(string(s))
}

type errorCmd struct {
	err error
}

func (e errorCmd) Perform(repl *Repl) {
	repl.MsgErr(e.err)
}

func quit(repl *Repl) {
	_, _ = fmt.Fprintf(repl.out(), "Done with %s", repl.CurrentCtx().Name)
	repl.FinishCtx()
}

func nop(_ *Repl) {}

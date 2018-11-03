package console

import (
	"fmt"
	"io"
)

// OutputInterface
type OutputInterface interface {
	Print(a ...interface{})
	Println(a ...interface{})
	Printf(format string, a ...interface{})
	Error(a ...interface{})
	Errorln(a ...interface{})
	Errorf(format string, a ...interface{})
}

// output
type output struct {
	out io.Writer
	err io.Writer
}

func (o output) Print(a ...interface{}) {
	o.write(o.out, a...)
}

func (o output) Println(a ...interface{}) {
	o.write(o.out, fmt.Sprintln(a...))
}

func (o output) Printf(format string, a ...interface{}) {
	o.write(o.out, fmt.Sprintf(format, a...))
}

func (o output) Error(a ...interface{}) {
	o.write(o.err, a...)
}

func (o output) Errorln(a ...interface{}) {
	o.write(o.err, fmt.Sprintln(a...))
}

func (o output) Errorf(format string, a ...interface{}) {
	o.write(o.err, fmt.Sprintf(format, a...))
}

func (o output) write(w io.Writer, a ...interface{}) {
	if _, err := fmt.Fprint(w, a...); err != nil {
		panic(fmt.Sprintf("gowl: cannot write to writer: %s", err.Error()))
	}
}

func NewOutput(out, err io.Writer) OutputInterface {
	return &output{out, err}
}

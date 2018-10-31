package gowl

import (
	"fmt"
	"io"
)

// CommandInterface
type CommandInterface interface {
	Name() string
	Help() string
	Execute(out OutputInterface)
}

// command
type command struct {
	name    string
	help    string
	handler func(out OutputInterface)
}

func (c *command) Name() string {
	return c.name
}

func (c *command) Help() string {
	return c.help
}

func (c *command) Execute(out OutputInterface) {
	c.handler(out)
}

func NewCommand(name, help string, handler func(out OutputInterface)) CommandInterface {
	return &command{name, help, handler}
}

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
	stdout io.Writer
	stderr io.Writer
}

func (o output) Print(a ...interface{}) {
	o.write(o.stdout, a...)
}

func (o output) Println(a ...interface{}) {
	o.write(o.stdout, fmt.Sprintln(a...))
}

func (o output) Printf(format string, a ...interface{}) {
	o.write(o.stdout, fmt.Sprintf(format, a...))
}

func (o output) Error(a ...interface{}) {
	o.write(o.stderr, a...)
}

func (o output) Errorln(a ...interface{}) {
	o.write(o.stderr, fmt.Sprintln(a...))
}

func (o output) Errorf(format string, a ...interface{}) {
	o.write(o.stderr, fmt.Sprintf(format, a...))
}

func (o output) write(w io.Writer, a ...interface{}) {
	if _, err := fmt.Fprint(w, a...); err != nil {
		panic(fmt.Sprintf("gowl: cannot write to writer: %s", err.Error()))
	}
}

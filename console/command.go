package console

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

package gowl

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/lokhman/gowl/console"
)

const (
	ServerName = "gowl/1.0"
)

const (
	Development = "dev"
	Staging     = "staging"
	Production  = "prod"
)

var kernel struct {
	servers  sync.Map
	commands sync.Map
}

var execPath string

var (
	_env    = flag.String("env", Production, "set executable environment")
	_debug  = flag.Bool("debug", false, "enable debug mode")
	_quiet  = flag.Bool("quiet", false, "quiet mode")
	_server = flag.String("server", "", "filter by server address")
)

var (
	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr
)

var (
	Debug = log.New(stdout, "[gowl] ", log.LstdFlags)
	Error = log.New(stderr, "[gowl] ", log.LstdFlags|log.Lshortfile)
)

func init() {
	var exec string
	if fn, err := os.Executable(); err == nil {
		execPath, exec = filepath.Split(fn)
	}

	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage: %s [flags] <command>", exec)
		fmt.Fprintln(out, "\n\nThe commands are:")
		commands := make([]string, 0)
		kernel.commands.Range(func(_, command interface{}) bool {
			c := command.(console.CommandInterface)
			help := strings.Replace(c.Help(), "\n", "\n    \t", -1)
			commands = append(commands, "  "+c.Name()+"\n    \t"+help)
			return true
		})
		sort.Strings(commands)
		fmt.Fprint(out, strings.Join(commands, "\n"))
		fmt.Fprintln(out, "\n\nThe flags are:")
		flag.PrintDefaults()
	}

	RegisterCommand(console.NewCommand("run", "run registered servers", runCommand))
	RegisterCommand(console.NewCommand("info", "display information about registered servers", infoCommand))
}

func ExecPath() string {
	return execPath
}

func EnvMode() string {
	return *_env
}

func DebugMode() bool {
	return *_debug
}

func RegisterCommand(command console.CommandInterface) {
	name := command.Name()
	if _, ok := kernel.commands.Load(name); ok {
		panic(fmt.Sprintf(`gowl: command "%s" is already registered`, name))
	}
	kernel.commands.Store(name, command)
}

func RegisterServer(server ServerInterface) {
	addr := server.Config().Addr
	if _, ok := kernel.servers.Load(addr); ok {
		panic(fmt.Sprintf(`gowl: server with address "%s" is already registered`, addr))
	}
	kernel.servers.Store(addr, server)
}

func Run(server ...ServerInterface) {
	flag.Parse()

	if *_quiet {
		stdout = ioutil.Discard
		stderr = ioutil.Discard
	}

	for _, server := range server {
		RegisterServer(server)
	}

	if name := flag.Arg(0); name != "" {
		out := console.NewOutput(stdout, stderr)
		getCommand(name).Execute(out)
		return
	}
	flag.Usage()
}

func getServer(addr string) ServerInterface {
	server, ok := kernel.servers.Load(addr)
	if !ok {
		fatal(`Server "%s" is not registered`, addr)
	}
	return server.(ServerInterface)
}

func getCommand(name string) console.CommandInterface {
	command, ok := kernel.commands.Load(name)
	if !ok {
		fatal(`Command "%s" is not registered`, name)
	}
	return command.(console.CommandInterface)
}

func fatal(format string, a ...interface{}) {
	fmt.Fprintln(stderr, fmt.Sprintf(format, a...))
	os.Exit(1)
}

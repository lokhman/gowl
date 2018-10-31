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

	"golang.org/x/sync/errgroup"
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
		str := "Usage: " + exec + " [flags] <command>\n\n"
		str += "The commands are:\n"
		commands := make([]string, 0)
		kernel.commands.Range(func(_, command interface{}) bool {
			c := command.(CommandInterface)
			help := strings.Replace(c.Help(), "\n", "\n    \t", -1)
			commands = append(commands, "  "+c.Name()+"\n    \t"+help)
			return true
		})
		sort.Strings(commands)
		str += strings.Join(commands, "\n") + "\n\n"
		str += "The flags are:\n"
		fmt.Fprint(flag.CommandLine.Output(), str)
		flag.PrintDefaults()
	}

	RegisterCommand(NewCommand("run", "run registered servers", runCommand))
	RegisterCommand(NewCommand("info", "display information about registered servers", infoCommand))
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

func RegisterCommand(command CommandInterface) {
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
		out := &output{stdout, stderr}
		getCommand(name).Execute(out)
		return
	}
	flag.Usage()
}

func runCommand(out OutputInterface) {
	var stack errgroup.Group
	if addr := *_server; addr != "" {
		out.Printf("Starting server... %s\n", addr)
		stack.Go(getServer(addr).Listen)
	} else {
		count := 0
		kernel.servers.Range(func(addr, server interface{}) bool {
			out.Printf("Starting server... %s\n", addr.(string))
			stack.Go(server.(ServerInterface).Listen)
			count++
			return true
		})
		if count == 0 {
			out.Errorln("No registered servers")
			return
		}
	}
	if err := stack.Wait(); err != nil {
		Error.Fatal(err)
	}
}

func infoCommand(out OutputInterface) {
	if addr := *_server; addr != "" {
		out.Println(getServer(addr))
		return
	}

	i := 1
	kernel.servers.Range(func(_, server interface{}) bool {
		str := fmt.Sprintf("Server #%d", i)
		if i > 1 {
			out.Println()
		}
		out.Println(str)
		out.Println(strings.Repeat("-", len(str)))
		out.Println(server.(ServerInterface).String())
		i++
		return true
	})
}

func getServer(addr string) ServerInterface {
	server, ok := kernel.servers.Load(addr)
	if !ok {
		fatal(`Server "%s" is not registered`, addr)
	}
	return server.(ServerInterface)
}

func getCommand(name string) CommandInterface {
	command, ok := kernel.commands.Load(name)
	if !ok {
		fatal(`Command "%s" is not registered`, name)
	}
	return command.(CommandInterface)
}

func fatal(format string, a ...interface{}) {
	fmt.Fprintln(stderr, fmt.Sprintf(format, a...))
	os.Exit(1)
}

func noPanicWrapper(fn func() error) error {
	defer func() { _ = recover() }()
	return fn()
}

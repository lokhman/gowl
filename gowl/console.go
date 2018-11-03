package gowl

import (
	"fmt"
	"strings"

	"github.com/lokhman/gowl/console"
	"golang.org/x/sync/errgroup"
)

func runCommand(out console.OutputInterface) {
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

func infoCommand(out console.OutputInterface) {
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

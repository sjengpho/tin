package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sjengpho/tin/grpc"
	"github.com/sjengpho/tin/tin"
)

const help = `
Usage: tin-server --FLAG VALUE

FLAG		DEFAULT			DESCRIPTION
port		8717			Server port
`

func main() {
	flag.Usage = func() { fmt.Fprintf(os.Stderr, help) }
	port := flag.Int("port", 8717, "The server port")
	flag.Parse()

	if strings.ToLower(flag.Arg(0)) == "help" {
		fmt.Println(help)
		return
	}

	s := grpc.NewServer(tin.DefaultConfig())
	s.ListenAndServe(*port)
}

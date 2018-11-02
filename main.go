package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

var (
	ipfsAddr = flag.String("api", "/ip4/127.0.0.1/tcp/5001", "ipfs API server")
)

func run() error {
	sh := shell.NewShell(*ipfsAddr)
	//sh := shell.NewLocalShell()
	cid, err := sh.Add(strings.NewReader("myfoobarmkm"))
	if err != nil {
		return err
	}
	fmt.Printf("added %s\n", cid)
	return nil
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
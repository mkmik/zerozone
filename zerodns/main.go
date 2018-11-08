package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/bitnami-labs/zerozone/zerodns/zerozone"

	_ "github.com/coredns/coredns/plugin/errors"
	_ "github.com/coredns/coredns/plugin/forward"
	_ "github.com/coredns/coredns/plugin/log"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
	shell "github.com/ipfs/go-ipfs-api"
)

var (
	ipfsAddr = flag.String("api", "/ip4/127.0.0.1/tcp/5001", "ipfs API server")
)

func init() {
	dnsserver.Directives = append([]string{"zerozone"}, dnsserver.Directives...)
}

func run() error {
	sh := shell.NewShell(*ipfsAddr)
	cid, err := sh.Add(strings.NewReader("myfoobarmkm"))
	if err != nil {
		return err
	}
	fmt.Printf("added %s\n", cid)

	fmt.Println("directives", dnsserver.Directives)
	coremain.Run()
	return nil
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
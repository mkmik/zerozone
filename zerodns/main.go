package main

import (
	"flag"
	"log"

	_ "github.com/bitnami-labs/zerozone/zerodns/zerozone"

	_ "github.com/coredns/coredns/plugin/errors"
	_ "github.com/coredns/coredns/plugin/forward"
	_ "github.com/coredns/coredns/plugin/log"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
)

func init() {
	dnsserver.Directives = append(dnsserver.Directives, "zerozone")
}

func run() error {
	coremain.Run()
	return nil
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
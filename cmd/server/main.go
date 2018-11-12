package main

import (
	"flag"
	"log"

	_ "github.com/bitnami-labs/zerozone/pkg/coredns/zerozone"

	_ "github.com/coredns/coredns/plugin/errors"
	_ "github.com/coredns/coredns/plugin/file"
	_ "github.com/coredns/coredns/plugin/forward"
	_ "github.com/coredns/coredns/plugin/log"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
)

func init() {
	dnsserver.Directives = []string{
		"metadata",
		"tls",
		"reload",
		"root",
		"debug",
		"trace",
		"health",
		"pprof",
		"prometheus",
		"errors",
		"log",
		"dnstap",

		"zerozone",
		"file",
	}
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

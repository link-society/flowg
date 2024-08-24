package main

import (
	"flag"
	"os"
)

var (
	bindAddress *string
	logDir      *string
	configDir   *string
	verbose     *bool
)

func init() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	bindAddress = flag.String(
		"bind",
		":5080",
		"Address to bind the server to",
	)

	logDir = flag.String(
		"log-dir",
		"./data/logs",
		"Path to the log database directory",
	)

	configDir = flag.String(
		"config-dir",
		"./data/config",
		"Path to the config directory",
	)

	verbose = flag.Bool(
		"verbose",
		false,
		"Enable verbose logging",
	)

	flag.Parse()
}

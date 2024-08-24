package main

import "flag"

var (
	bindAddress = flag.String("bind", ":5080", "Address to bind the server to")
	logDir      = flag.String("log-dir", "./data/logs", "Path to the log database directory")
	configDir   = flag.String("config-dir", "./data/config", "Path to the config directory")
	verbose     = flag.Bool("verbose", false, "Enable verbose logging")
)

func init() {
	flag.Parse()
}

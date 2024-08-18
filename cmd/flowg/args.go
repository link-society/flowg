package main

import "flag"

var (
	bindAddress = flag.String("bind", ":5080", "Address to bind the server to")
	dbPath      = flag.String("db", "./data/logs", "Path to the log database directory")
	configDir   = flag.String("config", "./data/config", "Path to the config directory")
)

func init() {
	flag.Parse()
}

package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// Version is the version of the application
	Version            = "0.0.1"
	ServerName  string = "consumer"
	flagconf    string
	ServerID, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", fmt.Sprintf("config path, eg: -conf %s.yaml", ServerName))
}

func main() {
	flag.Parse()
}

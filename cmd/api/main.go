package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yanglunara/im/pkg/common/cmd"
)

func main() {
	if err := cmd.NewApiCommand().Execute(); err != nil {
		progName := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "%s exit -1: %+v\n", progName, err)
		os.Exit(-1)
	}
}

package main

import (
	"os"

	"github.com/0ghny/gitconfig/internal/cmd"
)

var (
	version = "dev"
	// https://goreleaser.com/cookbooks/using-main.version
	// commit  = "none"
	// date    = "unknown"
	// builtBy = "unknown"
)

func main() {
	root := cmd.RootCmd(version)
	if err := root.Execute(); err != nil {
		//log.Fatal(err)
		//fmt.Println(err)
		os.Exit(1)
	}
}

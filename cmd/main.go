package main

import (
	"os"

	"github.com/0ghny/gitconfigs/internal/cmd"
)

var (
	version = "0.0.0"
)

func main() {
	root := cmd.RootCmd(version)
	if err := root.Execute(); err != nil {
		//log.Fatal(err)
		//fmt.Println(err)
		os.Exit(1)
	}
}

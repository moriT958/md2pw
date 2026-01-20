package main

import (
	"os"

	"github.com/moriT958/md2pw/internal/cli"
)

func main() {
	c := cli.New(os.Stdin, os.Stdout, os.Stderr)
	os.Exit(c.Run(os.Args))
}

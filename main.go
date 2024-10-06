package main

import (
	"fmt"
	"os"

	"github.com/cuctemeh/rstream-consumer/cmd"
)

func main() {
	err := cmd.NewRootCMD().Execute()
	if err == nil {
		return
	}

	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

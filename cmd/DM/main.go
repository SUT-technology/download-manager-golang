package main

import (
	"fmt"
	"os"

	"github.com/SUT-technology/download-manager-golang/cmd/DM/command"
)

func main() {
	err := command.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

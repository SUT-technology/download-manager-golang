package main

import (
	"fmt"
	"github.com/SUT-technology/download-manager-golang/cmd/DM/command"
	"os"
)

func main() {
	err := command.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

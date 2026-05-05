package main

import (
	"fmt"
	"os"
)

var version = "0.1.0"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("komyzi version", version)
		return
	}

	fmt.Println("Komyzi - AI Agent Configuration Manager")
	fmt.Println("Run 'komyzi --help' for usage information.")
}
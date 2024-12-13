package main

import (
	"fmt"
	"os"

	"go.datumapis.com/datum/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

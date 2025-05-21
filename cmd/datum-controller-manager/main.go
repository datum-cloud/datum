package main

import (
	"os"

	"k8s.io/component-base/cli"

	"go.datum.net/datum/cmd/datum-controller-manager/app"
)

func main() {
	command := app.NewCommand()
	code := cli.Run(command)
	os.Exit(code)
}

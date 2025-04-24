package main

import (
	"os"

	"k8s.io/component-base/cli"

	"go.datumapis.com/datum/cmd/datum-apiserver/app"
)

// Implementation approach follows work done in https://github.com/kubernetes/kubernetes/pull/126260
func main() {
	command := app.NewCommand()
	code := cli.Run(command)
	os.Exit(code)
}

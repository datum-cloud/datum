package main

import (
	"fmt"
	"os"

	"go.datum.net/datum/cmd/datum-authorization-webhook/app"
)

func main() {
	if err := app.NewWebhook().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

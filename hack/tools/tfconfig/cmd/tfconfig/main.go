package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/criticalstack/crit/hack/tools/tfconfig/cmd/tfconfig/app/create"
)

var cmd = &cobra.Command{
	Use:   "tfconfig",
	Short: "tfconfig is a tool for creating Terraform variable definitions files.",
}

func main() {
	cmd.AddCommand(create.NewCommand())
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

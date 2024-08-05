package main

import (
	"fmt"

	"github.com/oversize/txmon/pkg/txmon"
	"github.com/spf13/cobra"
)

func apiCommand() *cobra.Command {
	apiCommand := cobra.Command{
		Use:   "api",
		Short: "Run the api",
		Run: _apiCommand,
	}

	//apiCommand.AddCommand(
	//	apiSubCommand(),
	//)
	return &apiCommand
}

func _apiCommand(cmd *cobra.Command, args []string) {
	fmt.Print("Starting ")
	server := txmon.NewAPIServer(
		"localhost:9999",
	)
	server.Run()
}
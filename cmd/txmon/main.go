package main

import (
	"fmt"
	"os"
	"time"

	"github.com/oversize/txmon/pkg/txmon"
	"github.com/spf13/cobra"
)

const programName = "txmon"

func main() {

	rootCmd := &cobra.Command{
		Use: programName,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: false,
		},
		Run: rootCommand,
	}

	// Add some sub command
	rootCmd.AddCommand(
		apiCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func rootCommand(cmd *cobra.Command, args []string) {
	errorChan := make(chan error)
	go func() {
		for {
			err := <-errorChan
			fmt.Printf("ERROR: async: %s", err)
		}
	}()

	oConn, err := txmon.GetConnection(errorChan)
	if err != nil {
		fmt.Printf("failed to connect to node: %s", err)
	}

	for {
		// fmt.Println("Next")
		err = oConn.LocalTxMonitor().Client.Acquire()
		if err != nil {
			fmt.Println("failed to acquire mempool")
		}
		txmon.GetTransactions(oConn)
		err = oConn.LocalTxMonitor().Client.Release()
		if err != nil {
			fmt.Println("failed to release acquired mempool")
		}
		time.Sleep(time.Millisecond * 200)
	}
}

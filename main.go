package main

import (
	"fmt"

	ouroboros "github.com/blinklabs-io/gouroboros"
	"github.com/blinklabs-io/gouroboros/ledger"
)

func getConnection(errorChan chan error) (*ouroboros.Connection, error) {
	// Creates the Connection using oruouprbs new
	oConn, err := ouroboros.NewConnection(
		ouroboros.WithNetworkMagic(764824073),
		ouroboros.WithErrorChan(errorChan),
		ouroboros.WithNodeToNode(false),
		ouroboros.WithKeepAlive(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failure creating ouroboros connection: %s", err)
	}

	// dials the connection
	// err = oConn.Dial("tcp", "london1.preprod.relays.cardano.network:3001")
	err = oConn.Dial("unix", "/home/msch/src/cf/txmon/node.socket")
	if err != nil {
		return nil, fmt.Errorf(
			"failure connecting to node via TCP: %s",
			err,
		)
	}

	return oConn, nil
}

func getTransactions(oConn *ouroboros.Connection) {
	if oConn == nil {
		return
	}

	oConn.LocalTxMonitor().Client.Start()
	for {
		txRawBytes, err := oConn.LocalTxMonitor().Client.NextTx()
		if err != nil {
			fmt.Print(err.Error())
		}
		// size := len(txRawBytes)
		// fmt.Printf("Transaction Bytes %s", size)
		txType, err := ledger.DetermineTransactionType(txRawBytes)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		tx, err := ledger.NewTransactionFromCbor(txType, txRawBytes)
		if err != nil {
			fmt.Print(err.Error())
			return
		}

		fmt.Printf("Transactions: %s\n", tx.Hash())

	}
}

func main() {
	errorChan := make(chan error)
	go func() {
		for {
			err := <-errorChan
			fmt.Printf("ERROR: async: %s", err)
		}
	}()

	oConn, err := getConnection(errorChan)
	if err != nil {
		fmt.Printf("failed to connect to node: %s", err)
	}
	// time.Sleep(1)
	getTransactions(oConn)
}

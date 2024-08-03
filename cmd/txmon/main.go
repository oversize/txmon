package main

import (
	"fmt"
	"time"

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

	for {
		txRawBytes, err := oConn.LocalTxMonitor().Client.NextTx()
		if err != nil {
			fmt.Print(err.Error())
		}
		if len(txRawBytes) == 0 {
			return
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

	for {

		// fmt.Println("Next")
		err = oConn.LocalTxMonitor().Client.Acquire()
		if err != nil {
			fmt.Println("failed to acquire mempool")
		}
		getTransactions(oConn)
		err = oConn.LocalTxMonitor().Client.Release()
		if err != nil {
			fmt.Println("failed to release acquired mempool")
		}
		time.Sleep(time.Millisecond * 200)
	}

}

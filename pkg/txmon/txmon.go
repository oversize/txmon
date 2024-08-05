package txmon

import (
	"fmt"
	"os"

	ouroboros "github.com/blinklabs-io/gouroboros"
	"github.com/blinklabs-io/gouroboros/ledger"
)

func GetConnection(errorChan chan error) (*ouroboros.Connection, error) {
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

	// err = oConn.Dial("tcp", "london1.preprod.relays.cardano.network:3001")
	err = oConn.Dial("unix", os.Getenv("CARDANO_NODE_SOCKET_PATH"))
	if err != nil {
		return nil, fmt.Errorf(
			"failure connecting to node via TCP: %s",
			err,
		)
	}

	return oConn, nil
}

func GetTransactions(oConn *ouroboros.Connection) {
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

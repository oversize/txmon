# TX Monitoring

Simple draft to a tx-monitoring system for cardano. Based on the ideas
of blockperf. But at this point its more me learning go cardano and gouroboros.

```bash
# Create a local socket to a remote socket
ssh -L /home/msch/src/cf/txmon/node.socket:/opt/cardano/cnode/sockets/node.socket ubuntu@cardano.node
```

## Client

txmon can either run in client mode or in server mode. Client mode means, that
it starts to listen in the given node socket (txp or unix) and take snapshots
of the txs it sees in the local mempool. That is already mostly implemented
in cmd/txmon/main.go rootCommand() and then in the pkg/txmon. GetConection
connects to the local node and getTransactions takes that connection and
will poll(wait?) for tharnascions to be in the mempool.

Its currently two loops. One outer loop that acquires the lock to the
txmonitor which then blocks until there are (new) transactions in it. All of
 these will then be iterated until the current snapshot(?) is empty and
 a a new cycle begins (needing a new acquire of the lock).

 Every transaction it sees, it needs to take the tx id and the time it saw it
 first. So it needs to store which tx IDs it already saw. Just store that
 in memory for now. There is no need to persist this on the client side.
 Every time the lock is released, submit the transaction ids and their times
 to the central server.

 ### Local transaction store

 In memory map of all the transactions seen.

```go
type Transaction struct {
    txid string,
    appeared time.Time,  // as accurate as possible ...
                        // unixtimestamp would work but the more precision
                        // the better
    sent bool,      // Flag indicating the transaction has been submitted
}
txs := map[string]Transaction{
    "tx12345": {
        txid = "tx12345",
        appeared = now(),
        sent = false,
    },
    "tx6789a": {
        txid = "tx6789a",
        appeared = later(),
        sent = false,
    },
}
```
only add a given transaction if not alread inserted.
I will probably need to have a second structure that holds the txids that
are already sent to the server. I do not want to send a siingle txid twice
(from the same client!). Or just have the Transaction item in the list know
 whether  it was already setmt

Again, read up on how that txmonitring lock acquirering is supposed to work.
I believe i will not get around that process/flow. so i think its best
to always send the "current snapshot" once one such loop finishes. However,
a single loop might not produce new txids at all.

Once the inner loop finishes get the list of txids that zou need to send
and  create the list to hold them. Go through the whole txs array and
find the Transactions that are not yet sent.

```json
{
    relay_ip: "x.x.x.x", // IP Address of the client these came from
    node_version: "9.1.0", // cardano-node version of

    // Some notion of the "fullness" of the local mempool in that given cycle?

    transactions: [
        { txid: "12312312312312", appeard: "micro second datetime"}
        { txid: "12312312312312", appeard: "micro second datetime"}
        { txid: "12312312312312", appeard: "micro second datetime"}
    ]
}
```


## Server

The server receives tx samples from all clients over a simple http endpoint.
That is initially only one endpoint? One that receives above json. One such
json from every client roughly every few seconds.

The server then stores the above data like this (pseudo code)

### Tables

Does it make sens to store the samples as such? I want to be able to see
when any given transaction is seen first but also when other
relays then saw it. So that i can tell when a transaction was inserted into
the mempool but also determine when other relays saw it as well as the time
when it was finally adopted on chain.

**Transactions**

I need to know the transactions.
```sql
CREATE TABLE transactions (
id INT UNSIGNED NOT NULL AUTO_INCREMENT,
txid  "the tx id from chan"
createdAt "When it was created in server (not seen first)"
adopted "flag indicating it was actually inserted? Updated
        regluarly from dbsync or so?"

)
```

**Appearances**


```sql

```






```json
{

}
```

## Sampling Mempool Contents

TXMon connects to the local node and creates a sample of the current amempool
contents. It stores the transactions ids it sees and when it saw them first.



## Golang Resources

https://github.com/sikozonpc/go-rest-api


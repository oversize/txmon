# TX Monitoring

Simple draft to a tx-monitoring system for cardano. Based on the ideas
of blockperf. But at this point its more me learning go cardano and gouroboros.

```bash
# Create a local socket to a remote socket
ssh -L /home/msch/src/cf/txmon/node.socket:/opt/cardano/cnode/sockets/node.socket ubuntu@cardano.node
```

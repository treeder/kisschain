# KissChain - Keep it Simple Stupid Blockchain

## Sending commands

Create a transaction:

```sh
curl "localhost:8080/txion" \
     -H "Content-Type: application/json" \
     -d '{"from": "akjflw", "to":"fjlakdj", "amount": 3}'
```

Mine a block:

```sh
curl localhost:8080/mine
```

TRY SOME OF THESE TO TEST THE NETWORKING PART

Get blockchain

curl http://localhost:3001/blocks

Create block

curl -H "Content-type:application/json" --data '{"data" : "Some data to the first block"}' http://localhost:3001/mineBlock

Add peer

curl -H "Content-type:application/json" --data '{"peer" : "ws://localhost:6001"}' http://localhost:3001/addPeer

Query connected peers

curl http://localhost:3001/peers

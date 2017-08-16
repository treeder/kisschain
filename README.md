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

# bitcoin-stats

This is a bitcoin block reporting tool.  It allows you to query the current block height, show a particular block, and generate csv file for a range of blocks.  The block data includes the following fields:
```
type Block struct {
  hash string
  confirmations uint64
  strippedSize int32
  size int32
  weight int32
  height int64
  version int32
  merkleRoot string
  txnCount int
  reward float64
  txnFeeTotal float64
  totalOutputValue float64
  time time.Time
  nonce uint32
  bits string
  difficulty float64
  previousBlockHash string
  nextBlockHash string
}
```

This tool must be used in conjuction with [btcd](https://github.com/btcsuite/btcd).  It connects to local btcd server and obtains block data via websocket.  Run the following commands to launch the tool:
```
# Launch btcd with websocket enabled
btcd --rpcuser=[ANY_USER_NAME?] --rpcpass=[ANY_PASSWORD?]
# Wait for the btcd server to sync blocks from peers
# Go to bitcoin-stats/main.go, update the username and password there, then run:
go run bitcoin-stats/main.go
```

The tool will prompt you with options to select from.  Simply follow and explore.  A sample csv output is in the [output](output) folder.

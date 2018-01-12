package client

import (
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"io/ioutil"
	"path/filepath"
	"github.com/XiaoYang.Code4Fun/bitcoin-stats/data"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type Client struct {
  c *rpcclient.Client
}

func NewClient(user, pass string) (*Client, error) {
	// Connect to local btcd RPC server using websockets.
	btcdHomeDir := btcutil.AppDataDir("btcd", false)
	certs, err := ioutil.ReadFile(filepath.Join(btcdHomeDir, "rpc.cert"))
	if err != nil {
		return nil, err
	}
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8334",
		Endpoint:     "ws",
		User:         user,
		Pass:         pass,
		Certificates: certs,
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, err
	}
	return &Client{
		c: client,
	}, nil
}

func (c *Client) GetBlockCount() (int64, error) {
	blockCount, err := c.c.GetBlockCount()
	if err != nil {
		return 0, err
	}
	return blockCount, nil
}

func (c *Client) GetBlockByHeight(height int64) (*data.Block, error) {
  hash, err := c.c.GetBlockHash(height)
	if err != nil {
		return nil, err
	}
	return c.GetBlockByHash(hash.String())
}

func (c *Client) GetBlockByHash(hash string) (*data.Block, error) {
	h, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return nil, err
	}
  b, err := c.c.GetBlockVerboseTx(h)
	if err != nil {
		return nil, err
	}
	block, err := data.ParseBlock(b)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (c *Client) Close() {
	c.c.Shutdown()
}





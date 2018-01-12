package data

import (
	"time"
	"github.com/btcsuite/btcd/btcjson"
	"bytes"
	"fmt"
	"math"
)

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

func ParseBlock(b *btcjson.GetBlockVerboseResult) (*Block, error) {
  block := &Block{
		hash: b.Hash,
		confirmations: b.Confirmations,
		strippedSize: b.StrippedSize,
		size: b.Size,
		weight: b.Weight,
		height: b.Height,
		version: b.Version,
		merkleRoot: b.MerkleRoot,
		txnCount: len(b.RawTx),
		reward: 50 / math.Pow(2, float64((b.Height - 1) / 210000)),
		time: time.Unix(b.Time, 0),
		nonce: b.Nonce,
		bits: b.Bits,
		difficulty: b.Difficulty,
		previousBlockHash: b.PreviousHash,
		nextBlockHash: b.NextHash,
	}
	// Calculate txnFeeTotal, totalOutputValue
  for _, tx := range b.RawTx {
		for _, out := range tx.Vout {
			block.totalOutputValue += out.Value
		}
	}
	block.txnFeeTotal = block.totalOutputValue - block.reward
	return block, nil
}

func (b *Block) GetHeight() int64 {
	return b.height
}

func (b *Block) GetNextBlockHash() string {
	return b.nextBlockHash
}

func CSVHeader() string {
	return "hash,confirmations,strippedSize,size,weight,height,version,merkleRoot,txnCount,reward,txnFeeTotal,totalOutputValue,time,nonce,bits,difficulty,previousBlockHash,nextBlockHash\n"
}

func (b *Block) SingleBlockOutput() string {
  var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("hash: %v\n", b.hash))
	buffer.WriteString(fmt.Sprintf("confirmations: %v\n", b.confirmations))
	buffer.WriteString(fmt.Sprintf("hastrippedSizesh: %v\n", b.strippedSize))
	buffer.WriteString(fmt.Sprintf("size: %v\n", b.size))
	buffer.WriteString(fmt.Sprintf("weight: %v\n", b.weight))
	buffer.WriteString(fmt.Sprintf("height: %v\n", b.height))
	buffer.WriteString(fmt.Sprintf("version: %v\n", b.version))
	buffer.WriteString(fmt.Sprintf("merkleRoot: %v\n", b.merkleRoot))
	buffer.WriteString(fmt.Sprintf("txnCount: %v\n", b.txnCount))
	buffer.WriteString(fmt.Sprintf("reward: %v\n", b.reward))
	buffer.WriteString(fmt.Sprintf("txnFeeTotal: %v\n", b.txnFeeTotal))
	buffer.WriteString(fmt.Sprintf("totalOutputValue: %v\n", b.totalOutputValue))
	buffer.WriteString(fmt.Sprintf("time: %v\n", b.time))
	buffer.WriteString(fmt.Sprintf("nonce: %v\n", b.nonce))
	buffer.WriteString(fmt.Sprintf("bits: %v\n", b.bits))
	buffer.WriteString(fmt.Sprintf("difficulty: %v\n", b.difficulty))
	buffer.WriteString(fmt.Sprintf("previousBlockHash: %v\n", b.previousBlockHash))
	buffer.WriteString(fmt.Sprintf("nextBlockHash: %v\n", b.nextBlockHash))
	return buffer.String()
}

func (b *Block) CSVBlockOutput() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
		b.hash,
		b.confirmations,
		b.strippedSize,
		b.size,
		b.weight,
		b.height,
		b.version,
		b.merkleRoot,
		b.txnCount,
		b.reward,
		b.txnFeeTotal,
		b.totalOutputValue,
		b.time,
		b.nonce,
		b.bits,
		b.difficulty,
		b.previousBlockHash,
		b.nextBlockHash)
}

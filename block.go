package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

type Block struct {
	Index        uint64
	Timestamp    time.Time
	Data         []byte
	PreviousHash []byte
	Hash         []byte
}

func NewBlock(index uint64, timestamp time.Time, data []byte, prevHash []byte) *Block {
	b := &Block{
		Index:        index,
		Timestamp:    timestamp,
		Data:         data,
		PreviousHash: prevHash,
	}
	b.Hash = HashBlock(b)
	return b
}

func (b *Block) HashHex() string {
	return hex.EncodeToString(b.Hash)
}

func HashBlock(b *Block) []byte {
	h := sha256.New()
	h.Write([]byte(strconv.FormatUint(b.Index, 10)))
	h.Write([]byte(strconv.FormatInt(b.Timestamp.UnixNano(), 10)))
	h.Write(b.Data)
	h.Write(b.PreviousHash)

	// sha.update(str(self.index) +
	//            str(self.timestamp) +
	//            str(self.data) +
	//            str(self.previous_hash))
	// return sha.hexdigest()
	return h.Sum(nil)
}

func CreateGenesisBlock() *Block {
	return NewBlock(0, time.Now(), []byte("Genesis Block"), []byte("0"))
}

func NextBlock(prev *Block) *Block {
	i := prev.Index + 1
	b := NewBlock(
		i,
		time.Now(),
		[]byte("Hey! I'm block "+strconv.FormatUint(i, 10)),
		prev.Hash,
	)
	return b
}

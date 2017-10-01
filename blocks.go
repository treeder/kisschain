package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func blocksAll(w http.ResponseWriter, r *http.Request) {
	ret := map[string]interface{}{}
	ret["blockchain"] = blockchain
	writeJSON(w, http.StatusOK, ret)
}

func blocksLatest(w http.ResponseWriter, r *http.Request) {
	ret := map[string]interface{}{}
	ret["block"] = []*Block{getLatestBlock()}
	writeJSON(w, http.StatusOK, ret)
}

func getLatestBlock() *Block {
	return blockchain[len(blockchain)-1]
}

func blockReceived(w http.ResponseWriter, r *http.Request) {
	blockReceived := &BlockWrapper{}
	err := parseJSON(w, r, blockReceived)
	if err != nil {
		m := fmt.Sprintf("couldn't parse chain from peer %v - %v", r.RemoteAddr, err)
		log.Print(m)
		writeError(w, http.StatusBadRequest, errors.New(m))
		return
	}
	// If we allow receiving chains, need to sort by index
	ourLatestBlock := getLatestBlock()
	if blockReceived.Block.Index > ourLatestBlock.Index {
		// TODO: fill in missing blocks from peers
		if bytes.Compare(ourLatestBlock.Hash, blockReceived.Block.PreviousHash) == 0 {
			log.Println("we can append the received block to our chain")
			// TODO: need to sync these appends
			blockchain = append(blockchain, blockReceived.Block)
		} else {

		}
	} else {
		log.Println("received block is older than our current chain, do nothing.")
	}

	writeMessage(w, http.StatusOK, "pong")
}

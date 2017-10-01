package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/go-chi/chi"
)

var (
	// This node's blockchain copy
	blockchain []*Block
	// Store the transactions that this node has in a list
	thisNodesTransactions []*Transaction
	// A completely random address of the owner of this node
	minerAddress = "q3nf394hjg-random-miner-address-34nf3i4nflkn3oi"
	// Store the url data of every other node in the network so that we can communicate with them
	peerNodes []*Node
	// A variable to deciding if we're mining or not
	mining = true
)

func main() {
	blockchain = []*Block{CreateGenesisBlock()}

	r := chi.NewRouter()
	r.Get("/mine", mine)
	r.Post("/txion", transaction)
	r.Get("/blocks/latest", blocksLatest)
	r.Get("/blocks/all", blocksAll)
	r.Post("/blocks", blockReceived)

	r.Get("/peers", peers)
	r.Post("/peers/add", addPeer)
	r.Get("/ping", ping)
	port := ":" + getenv("PORT", "8080")
	fmt.Println("Starting KissChain on port", port)
	http.ListenAndServe(port, r)
}

func ping(w http.ResponseWriter, r *http.Request) {
	writeMessage(w, 200, "pong")
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func mine(w http.ResponseWriter, r *http.Request) {
	// Get the last proof of work
	lastBlock := blockchain[len(blockchain)-1]
	lastData := &POWData{}
	fmt.Println("lastData:", string(lastBlock.Data))
	err := json.Unmarshal(lastBlock.Data, lastData)
	if err != nil {
		log.Println("error parsing last data:", err)
		return
	}
	//   # Find the proof of work for
	//   # the current block being mined
	//   # Note: The program will hang here until a new
	//   #       proof of work is found
	proof := ProofOfWork(lastData.ProofOfWork)
	//   # Once we find a valid proof of work,
	//   # we know we can mine a block so
	//   # we reward the miner by adding a transaction
	thisNodesTransactions = append(thisNodesTransactions, &Transaction{
		From:   "network",
		To:     minerAddress,
		Amount: 1,
	})
	//   # Now we can gather the data needed
	//   # to create the new block
	newData := &POWData{
		ProofOfWork:  proof,
		Transactions: thisNodesTransactions,
	}
	newDataBytes, err := json.Marshal(newData)
	if err != nil {
		panic("error marshalling new data yo!")
	}
	minedBlock := NewBlock(
		lastBlock.Index+1,
		time.Now(),
		newDataBytes,
		lastBlock.Hash,
	)
	//   # Empty transaction list
	//   this_nodes_transactions[:] = []
	blockchain = append(blockchain, minedBlock)
	//   # Let the client know we mined a block
	writeJSON(w, 201, map[string]interface{}{
		"block": map[string]interface{}{
			"index":     minedBlock.Index,
			"timestamp": minedBlock.Timestamp.Format(time.RFC3339),
			"data":      minedBlock.Data,
			"hash":      minedBlock.HashHex(),
		},
		"message": "New block mined",
	})
}

func transaction(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t Transaction
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	thisNodesTransactions = append(thisNodesTransactions, &t)
	log.Printf("new transaction: %+v", t)

	writeJSON(w, 201, map[string]interface{}{
		"message": "Transaction successful",
	})
}

func findNewChains() ([]*Blockchain, error) {
	//  # Get the blockchains of every other node
	otherChains := []*Blockchain{}
	for _, p := range peerNodes {
		// # Get their chains using a GET request
		resp, err := http.Get(p.URL.String())
		if err != nil {
			log.Printf("error getting peers chains %v - %v", p.URL, err)
			return nil, err
		}
		defer resp.Body.Close()
		blockchain := &Blockchain{}
		err = parseJSONReader(resp.Body, blockchain)
		if err != nil {
			log.Printf("couldn't parse chain from peer %v - %v", p.URL, err)
			return nil, err
		}
		otherChains = append(otherChains, blockchain)
	}
	return otherChains, nil
}

func consensus() {
	//   # Get the blocks from other nodes
	otherChains, err := findNewChains()
	if err != nil {
		log.Println("Error getting other chains:", err)
		return
	}
	//   # If our chain isn't longest, then we store the longest chain
	longestChain := blockchain
	for _, chain := range otherChains {
		if len(longestChain) < len(chain.Blockchain) {
			longestChain = chain.Blockchain
		}
	}
	//   # If the longest chain isn't ours,then we stop mining and set
	//   # our chain to the longest one
	blockchain = longestChain
}

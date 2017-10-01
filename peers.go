package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Node struct {
	URL url.URL
}

func peers(w http.ResponseWriter, r *http.Request) {
	ret := map[string]interface{}{}
	ret["peers"] = peerNodes
	writeJSON(w, http.StatusOK, ret)
}

func addPeer(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t Node
	err := decoder.Decode(&t)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid input to add peer: %v", err))
		return
	}
	err = connectToPeers([]*Node{&t})
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("could not connect to this peer: %v", err))
		return
	}
	writeObject(w, 200, BasicResponse{
		Message: "Peer added successful",
	})
}

func connectToPeers(newPeers []*Node) error {
	for _, p := range newPeers {
		// test to ensure peer is OK
		resp, err := http.Get(fmt.Sprintf("%v/ping"))
		if err != nil {
			return err
		}
		br := &BasicResponse{}
		err = parseJSONReader(resp.Body, br)
		if err != nil {
			return fmt.Errorf("couldn't parse peers response: %v", err)
		}
		peerNodes = append(peerNodes, p)
		log.Printf("new peer added: %+v", p)
	}
	return nil
}

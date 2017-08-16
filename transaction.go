package main

type Transaction struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

// POWData is proof of work struct for the chain
type POWData struct {
	ProofOfWork  int `json:"proof_of_work"`
	Transactions []*Transaction
}

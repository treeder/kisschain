package main

func ProofOfWork(lastProof int) int {
	// Create a variable that we will use to find
	// our next proof of work
	incrementor := lastProof + 1
	// Keep incrementing the incrementor until
	// it's equal to a number divisible by 9
	// and the proof of work of the previous
	// block in the chain
	for {
		if incrementor%9 == 0 && (lastProof == 0 || incrementor%lastProof == 0) {
			break
		}
		incrementor++
	}
	// Once that number is found, we can return it as a proof of work
	return incrementor
}

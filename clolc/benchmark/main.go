package main

import (
	"flag"
)

func main() {
	benchPhasePtr := flag.String("phase", "", "Benchmark of CLOLC four phases: in, tr, ce, and rv")
	benchProcessPtr := flag.String("process", "", "submit_tx, read_tx, read_all_tx")
	numOrgPtr := flag.Int("numOrg", 2, "Number of organizations")
	numIterPtr := flag.Int("numIter", 10, "Number of iterations")
	numTXsPtr := flag.Int("numTXs", 100, "Number of transactions")
	flag.Parse()

	switch *benchPhasePtr {
	case "in":
		benchInitializeEpoch(*numOrgPtr, *numIterPtr)
	case "tr":
		switch *benchProcessPtr {
		case "submit_tx":
			benchTransactionRecordSubmitTX(*numTXsPtr, *numIterPtr)
		case "read_tx":
			benchTransactionRecordReadTX(*numTXsPtr, *numIterPtr)
		case "read_all_tx":
			benchTransactionRecordReadAllTXs(*numTXsPtr, *numIterPtr)
		}
	}
}

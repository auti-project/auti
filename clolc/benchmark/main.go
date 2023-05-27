package main

import (
	"flag"

	bf "github.com/auti-project/auti/clolc/benchmark/internal/benchfeature"
)

func main() {
	benchPhasePtr := flag.String("phase", "", "Benchmark of CLOLC four phases: in, tr, ce, and rv")
	benchProcessPtr := flag.String("process", "",
		"[tr]: submit_tx, read_tx, read_all_tx, commit_tx, accumulate, prepare_tx",
	)
	numOrgPtr := flag.Int("numOrg", 2, "Number of organizations")
	numIterPtr := flag.Int("numIter", 10, "Number of iterations")
	numTXsPtr := flag.Int("numTXs", 100, "Number of transactions")
	flag.Parse()

	switch *benchPhasePtr {
	case "in":
		bf.InitializeEpoch(*numOrgPtr, *numIterPtr)
	case "tr":
		switch *benchProcessPtr {
		case "submit_tx":
			if err := bf.TransactionRecordSubmitTX(*numTXsPtr, *numIterPtr); err != nil {
				return
			}
		case "prepare_tx":
			if err := bf.PrepareTX(*numTXsPtr); err != nil {
				return
			}
		case "read_tx":
			if err := bf.TransactionRecordReadTX(*numTXsPtr, *numIterPtr); err != nil {
				return
			}
		case "read_all_tx":
			if err := bf.TransactionRecordReadAllTXs(*numTXsPtr, *numIterPtr); err != nil {
				return
			}
		case "commit_tx":
			if err := bf.TransactionRecordCommitment(*numTXsPtr, *numIterPtr); err != nil {
				return
			}
		case "accumulate":
			if err := bf.TransactionRecordAccumulate(*numTXsPtr, *numIterPtr); err != nil {
				return
			}
		}
	}
}

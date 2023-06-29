package main

import (
	"flag"

	bf "github.com/auti-project/auti/benchmark/closc/internal/benchfeature"
	. "github.com/auti-project/auti/benchmark/closc/internal/flags"
)

func main() {
	benchPhasePtr := flag.String("phase", "", GetPhases())
	benchProcessPtr := flag.String("process", "", GetPhasesAndProcesses())
	numOrgPtr := flag.Int("numOrg", 2, "Number of organizations")
	numIterPtr := flag.Int("numIter", 10, "Number of iterations")
	numPtr := flag.Int("num", 100, "Number/Quantity/Depth")
	flag.Parse()

	var err error
	switch *benchPhasePtr {
	case PhaseInitialization:
		switch *benchProcessPtr {
		case ProcessINDefault:
			err = bf.InitializeEpoch(*numOrgPtr, *numIterPtr)
		case ProcessINRandGen:
			err = bf.InitializeRandGen(*numPtr, *numIterPtr)
		}
	case PhaseTransactionRecord:
		switch *benchProcessPtr {
		case ProcessTRCommitment:
			err = bf.TransactionRecordCommitment(*numPtr, *numIterPtr)
		case ProcessTRMerkleProofGen:
			err = bf.TransactionRecordMerkleProofGen(*numPtr, *numIterPtr)
		}
	}
	if err != nil {
		return
	}
}

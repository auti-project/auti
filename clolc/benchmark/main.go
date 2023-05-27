package main

import (
	"flag"

	bf "github.com/auti-project/auti/clolc/benchmark/internal/benchfeature"
)

func main() {
	benchPhasePtr := flag.String("phase", "", "Benchmark of CLOLC four phases: in, tr, ce, and rv")
	benchProcessPtr := flag.String("process", "",
		"[tr]: local_submit, local_read, local_read_all, local_prepare, commit_tx, accumulate, org_submit, org_read, org_read_all, org_prepare",
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
		case "local_submit":
			if err := bf.TransactionRecordLocalSubmitTX(*numTXsPtr, *numIterPtr); err != nil {
				return
			}
		case "local_prepare":
			if err := bf.PrepareLocalTX(*numTXsPtr); err != nil {
				return
			}
		case "local_read":
			if err := bf.TransactionRecordLocalReadTX(*numTXsPtr, *numIterPtr); err != nil {
				return
			}
		case "local_read_all":
			if err := bf.TransactionRecordLocalReadAllTXs(*numTXsPtr, *numIterPtr); err != nil {
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
		case "org_prepare":
			if err := bf.PrepareOrgTX(*numTXsPtr); err != nil {
				return
			}
		case "org_submit":
			if err := bf.TransactionRecordOrgSubmitTX(*numTXsPtr, *numIterPtr); err != nil {
				return
			}
		case "org_read":
			if err := bf.TransactionRecordOrgReadTX(*numTXsPtr, *numIterPtr); err != nil {
				return
			}
		case "org_read_all":
			if err := bf.TransactionRecordOrgReadAllTXs(*numTXsPtr, *numIterPtr); err != nil {
				return
			}
		}
	}
}

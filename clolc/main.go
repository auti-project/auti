package main

import (
	"flag"
	"log"

	bf "github.com/auti-project/auti/clolc/benchmark/internal/benchfeature"
)

func main() {
	benchPhasePtr := flag.String("phase", "", "in, tr, ce, rv")
	benchProcessPtr := flag.String("process", "",
		"[tr]: local_submit, local_read, local_read_all, local_prepare, commit_tx, accumulate, org_submit, org_read, org_read_all, org_prepare"+
			"[ce]: acc_commit, cal_b, cal_c, cal_d, encrypt, submit_tx, read_tx, read_all_txs, prepare_tx",
	)
	numOrgPtr := flag.Int("numOrg", 2, "Number of organizations")
	numIterPtr := flag.Int("numIter", 10, "Number of iterations")
	numTXsPtr := flag.Int("numTXs", 100, "Number of transactions")
	flag.Parse()

	var err error
	switch *benchPhasePtr {
	case "in":
		err = bf.InitializeEpoch(*numOrgPtr, *numIterPtr)
	case "tr":
		switch *benchProcessPtr {
		case "local_submit":
			err = bf.TransactionRecordLocalSubmitTX(*numTXsPtr, *numIterPtr)
		case "local_prepare":
			err = bf.PrepareLocalTX(*numTXsPtr)
		case "local_read":
			err = bf.TransactionRecordLocalReadTX(*numTXsPtr, *numIterPtr)
		case "local_read_all":
			err = bf.TransactionRecordLocalReadAllTXs(*numTXsPtr, *numIterPtr)
		case "commit_tx":
			err = bf.TransactionRecordCommitment(*numTXsPtr, *numIterPtr)
		case "accumulate":
			err = bf.TransactionRecordAccumulate(*numTXsPtr, *numIterPtr)
		case "org_prepare":
			err = bf.PrepareOrgTX(*numTXsPtr)
		case "org_submit":
			err = bf.TransactionRecordOrgSubmitTX(*numTXsPtr, *numIterPtr)
		case "org_read":
			err = bf.TransactionRecordOrgReadTX(*numTXsPtr, *numIterPtr)
		case "org_read_all":
			err = bf.TransactionRecordOrgReadAllTXs(*numTXsPtr, *numIterPtr)
		}
	case "ce":
		switch *benchProcessPtr {
		case "acc_commit":
			err = bf.ConsistencyExaminationAccumulateCommitment(*numOrgPtr, *numIterPtr)
		case "cal_b":
			err = bf.ConsistencyExaminationComputeB(*numOrgPtr, *numIterPtr)
		case "cal_c":
			err = bf.ConsistencyExaminationComputeC(*numOrgPtr, *numIterPtr)
		case "cal_d":
			err = bf.ConsistencyExaminationComputeD(*numOrgPtr, *numIterPtr)
		case "encrypt":
			err = bf.ConsistencyExaminationEncrypt(*numOrgPtr, *numIterPtr)
		case "submit_tx":
			err = bf.ConsistencyExaminationAudSubmitTX(*numTXsPtr, *numIterPtr)
		case "read_tx":
			err = bf.ConsistencyExaminationAudReadTX(*numTXsPtr, *numIterPtr)
		case "read_all_txs":
			err = bf.ConsistencyExaminationAudReadAllTXs(*numTXsPtr, *numIterPtr)
		case "prepare_tx":
			err = bf.PrepareAudTX(*numTXsPtr)
		}
	}
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

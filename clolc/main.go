package main

import (
	"flag"
	"log"

	bf "github.com/auti-project/auti/clolc/internal/benchfeature"
	. "github.com/auti-project/auti/clolc/internal/flag"
)

func main() {
	benchPhasePtr := flag.String("phase", "", GetPhases())
	benchProcessPtr := flag.String("process", "", GetPhasesAndProcesses())
	numOrgPtr := flag.Int("numOrg", 2, "Number of organizations")
	numIterPtr := flag.Int("numIter", 10, "Number of iterations")
	numTXsPtr := flag.Int("numTXs", 100, "Number of transactions")
	flag.Parse()

	var err error
	switch *benchPhasePtr {
	case PhaseInitialization:
		err = bf.InitializeEpoch(*numOrgPtr, *numIterPtr)
	case PhaseTransactionRecord:
		switch *benchProcessPtr {
		case ProcessTRLocalChainSubmit:
			err = bf.TransactionRecordLocalSubmitTX(*numTXsPtr, *numIterPtr)
		case ProcessTRLocalChainPrepare:
			err = bf.PrepareLocalTX(*numTXsPtr)
		case ProcessTRLocalChainRead:
			err = bf.TransactionRecordLocalReadTX(*numTXsPtr, *numIterPtr)
		case ProcessTRLocalChainReadAll:
			err = bf.TransactionRecordLocalReadAllTXs(*numTXsPtr, *numIterPtr)
		case ProcessTROrgChainSubmit:
			err = bf.TransactionRecordOrgSubmitTX(*numTXsPtr, *numIterPtr)
		case ProcessTROrgChainPrepare:
			err = bf.PrepareOrgTX(*numTXsPtr)
		case ProcessTROrgChainRead:
			err = bf.TransactionRecordOrgReadTX(*numTXsPtr, *numIterPtr)
		case ProcessTROrgChainReadAll:
			err = bf.TransactionRecordOrgReadAllTXs(*numTXsPtr, *numIterPtr)
		case ProcessTRCommitment:
			err = bf.TransactionRecordCommitment(*numTXsPtr, *numIterPtr)
		case ProcessTRAccumulate:
			err = bf.TransactionRecordAccumulate(*numTXsPtr, *numIterPtr)
		}
	case PhaseConsistencyExamination:
		switch *benchProcessPtr {
		case ProcessCEAccumulateCommitment:
			err = bf.ConsistencyExaminationAccumulateCommitment(*numOrgPtr, *numIterPtr)
		case ProcessCEComputeB:
			err = bf.ConsistencyExaminationComputeB(*numOrgPtr, *numIterPtr)
		case ProcessCEComputeC:
			err = bf.ConsistencyExaminationComputeC(*numOrgPtr, *numIterPtr)
		case ProcessCEComputeD:
			err = bf.ConsistencyExaminationComputeD(*numOrgPtr, *numIterPtr)
		case ProcessCEEncrypt:
			err = bf.ConsistencyExaminationEncrypt(*numOrgPtr, *numIterPtr)
		case ProcessCEAudChainSubmit:
			err = bf.ConsistencyExaminationAudSubmitTX(*numTXsPtr, *numIterPtr)
		case ProcessCEAudChainPrepare:
			err = bf.PrepareAudTX(*numTXsPtr)
		case ProcessCEAudChainRead:
			err = bf.ConsistencyExaminationAudReadTX(*numTXsPtr, *numIterPtr)
		case ProcessCEAudChainReadAll:
			err = bf.ConsistencyExaminationAudReadAllTXs(*numTXsPtr, *numIterPtr)
		case ProcessCECheck:
			err = bf.ConsistencyExaminationCheck(*numIterPtr)
		}
	default:
		log.Fatalf("Error: %v", "Invalid phase")
	}
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

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
		case ProcessTRLocalChainSubmit:
			err = bf.TransactionRecordLocalSubmitTX(*numPtr, *numIterPtr)
		case ProcessTRLocalChainPrepare:
			err = bf.TransactionRecordLocalPrepareTX(*numPtr)
		case ProcessTRLocalChainRead:
			err = bf.TransactionRecordLocalReadTX(*numPtr, *numIterPtr)
		case ProcessTRLocalChainReadAll:
			err = bf.TransactionRecordLocalReadAllTXs(*numPtr, *numIterPtr)
		case ProcessTRLocalChainCommitmentSubmit:
			err = bf.TransactionRecordLocalCommitmentSubmitTX(*numPtr, *numIterPtr)
		case ProcessTRLocalChainCommitmentPrepare:
			err = bf.TransactionRecordLocalCommitmentPrepareTX(*numPtr)
		case ProcessTRLocalCHainCommitmentRead:
			err = bf.TransactionRecordLocalCommitmentReadTX(*numPtr, *numIterPtr)
		case ProcessTRLocalChainCommitmentReadAll:
			err = bf.TransactionRecordLocalCommitmentReadAllTXs(*numPtr, *numIterPtr)
		case ProcessTROrgChainSubmit:
			err = bf.TransactionRecordOrgSubmitTX(*numPtr, *numIterPtr)
		case ProcessTROrgChainPrepare:
			err = bf.TransactionRecordOrgPrepareTX(*numPtr)
		case ProcessTROrgChainRead:
			err = bf.TransactionRecordOrgReadTX(*numPtr, *numIterPtr)
		case ProcessTROrgChainReadAll:
			err = bf.TransactionRecordOrgReadAllTXs(*numPtr, *numIterPtr)
		}
	case PhaseConsistencyExamination:
		switch *benchProcessPtr {
		case ProcessCEMerkleProofVerify:
			err = bf.ConsistencyExaminationMerkleProofVerify(*numPtr, *numIterPtr)
		case ProcessCEMerkleProofMerge:
			err = bf.ConsistencyExaminationMerkleProofMerge(*numPtr, *numIterPtr)
		case ProcessCESummarizeMerkleProofVerificationResults:
			err = bf.ConsistencyExaminationSummarizeMerkleProofVerificationResults(*numPtr, *numIterPtr)
		case ProcessCEVerifyCommitments:
			err = bf.ConsistencyExaminationVerifyCommitments(*numPtr, *numIterPtr)
		}
	}
	if err != nil {
		return
	}
}

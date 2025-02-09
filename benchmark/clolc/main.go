package main

import (
	"flag"
	"log"

	. "github.com/auti-project/auti/benchmark/clolc/internal/flag"
	"github.com/auti-project/auti/benchmark/clolc/internal/task"
)

func main() {
	benchPhasePtr := flag.String("phase", "", GetPhases())
	benchProcessPtr := flag.String("process", "", GetPhasesAndProcesses())
	numOrgPtr := flag.Int("numOrg", 2, "Number of organizations")
	numIterPtr := flag.Int("numIter", 10, "Number of iterations")
	numTXsPtr := flag.Int("numTXs", 100, "Number of transactions")
	numRoutinesPtr := flag.Int("numRoutines", 0, "Number of routines")
	flag.Parse()

	var err error
	switch *benchPhasePtr {
	case PhaseInitialization:
		switch *benchProcessPtr {
		case ProcessINDefault:
			err = task.INDefault(*numOrgPtr, *numIterPtr)
		}
	case PhaseTransactionRecord:
		switch *benchProcessPtr {
		case ProcessTRLocalChainSubmit:
			err = task.TRLocalSubmitTX(*numTXsPtr, *numIterPtr)
		case ProcessTRLocalChainPrepare:
			err = task.TRLocalPrepareTX(*numTXsPtr)
		case ProcessTRLocalChainRead:
			err = task.TRLocalReadTX(*numTXsPtr, *numIterPtr)
		case ProcessTRLocalChainReadAll:
			err = task.TRLocalReadAllTXs(*numTXsPtr, *numIterPtr)
		case ProcessTROrgChainSubmit:
			err = task.TROrgSubmitTX(*numTXsPtr, *numIterPtr)
		case ProcessTROrgChainPrepare:
			err = task.TROrgPrepareTX(*numTXsPtr)
		case ProcessTROrgChainRead:
			err = task.TROrgReadTX(*numTXsPtr, *numIterPtr)
		case ProcessTROrgChainReadAll:
			err = task.TROrgReadAllTXs(*numTXsPtr, *numIterPtr)
		case ProcessTRCommitment:
			err = task.TRCommitment(*numTXsPtr, *numIterPtr)
		case ProcessTRAccumulate:
			err = task.TRAccumulate(*numTXsPtr, *numIterPtr)
		}
	case PhaseConsistencyExamination:
		switch *benchProcessPtr {
		case ProcessCEAccumulateCommitment:
			err = task.CEAccumulateCommitment(*numOrgPtr, *numIterPtr)
		case ProcessCEComputeB:
			err = task.CEComputeB(*numOrgPtr, *numIterPtr)
		case ProcessCEComputeC:
			err = task.CEComputeC(*numOrgPtr, *numIterPtr)
		case ProcessCEComputeD:
			err = task.CEComputeD(*numOrgPtr, *numIterPtr)
		case ProcessCEEncrypt:
			err = task.CEEncrypt(*numOrgPtr, *numIterPtr)
		case ProcessCEDecrypt:
			err = task.CEDecrypt(*numIterPtr)
		case ProcessCEAudChainSubmit:
			err = task.CEAudSubmitTX(*numTXsPtr, *numIterPtr)
		case ProcessCEAudChainPrepare:
			err = task.CEAudPrepareTX(*numTXsPtr)
		case ProcessCEAudChainRead:
			err = task.CEAudReadTX(*numTXsPtr, *numIterPtr)
		case ProcessCEAudChainReadAll:
			err = task.CEAudReadAllTXs(*numTXsPtr, *numIterPtr)
		case ProcessCECheck:
			err = task.CECheck(*numIterPtr)
		case ProcessCEConsistencyExaminationPartOneParallel:
			err = task.CEBatchConsistencyExaminationPartOne(*numIterPtr, *numRoutinesPtr)
		case ProcessCEConsistencyExaminationPartTwoParallel:
			err = task.CEBatchConsistencyExaminationPartTwo(*numIterPtr, *numRoutinesPtr)
		}
	case PhaseResultVerification:
		switch *benchProcessPtr {
		case ProcessRVVerifyOrgAndAudResult:
			err = task.RVVerifyOrgAndAudResult(*numOrgPtr, *numIterPtr)
		case ProcessRVVerifyAuditPairResult:
			err = task.RVVerifyAuditPairResult(*numOrgPtr, *numIterPtr)
		case ProcessRVDecryptParallel:
			err = task.RVBatchDecrypt(*numIterPtr, *numRoutinesPtr)
		case ProcessRVCheckOrgAndAudPairParallel:
			err = task.RVBatchCheckOrgAndAudPair(*numIterPtr, *numRoutinesPtr)
		case ProcessRVCheckAudPairParallel:
			err = task.RVBatchCheckAudPair(*numIterPtr, *numRoutinesPtr)
		}

	default:
		log.Fatalf("Error: %v", "Invalid phase")
	}
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

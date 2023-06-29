package flags

const (
	PhaseInitialization         = "in"
	PhaseTransactionRecord      = "tr"
	PhaseConsistencyExamination = "ce"
	PhaseResultVerification     = "rv"

	ProcessINDefault = "default"
	ProcessINRandGen = "rand_gen"

	ProcessTRCommitment     = "commitment"
	ProcessTRMerkleProofGen = "merkle_proof_gen"
)

var PhaseProcessMap = map[string][]string{
	PhaseInitialization: {
		ProcessINDefault,
		ProcessINRandGen,
	},
	PhaseTransactionRecord: {
		ProcessTRCommitment,
		ProcessTRMerkleProofGen,
	},
}

func GetPhases() string {
	var phases string
	for key := range PhaseProcessMap {
		phases += key + ", "
	}
	return phases
}

func GetPhasesAndProcesses() string {
	var result string
	for key, val := range PhaseProcessMap {
		result += key + ": "
		for _, process := range val {
			result += process + ", "
		}
		result += "\n"
	}
	// remove the last "\n"
	result = result[:len(result)-1]
	return result
}

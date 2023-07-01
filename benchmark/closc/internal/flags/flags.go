package flags

const (
	PhaseInitialization         = "in"
	PhaseTransactionRecord      = "tr"
	PhaseConsistencyExamination = "ce"
	PhaseResultVerification     = "rv"

	ProcessINDefault = "default"
	ProcessINRandGen = "rand_gen"

	ProcessTRCommitment                  = "commitment"
	ProcessTRMerkleProofGen              = "merkle_proof_gen"
	ProcessTRLocalChainSubmit            = "local_submit"
	ProcessTRLocalChainPrepare           = "local_prepare"
	ProcessTRLocalChainRead              = "local_read"
	ProcessTRLocalChainReadAll           = "local_read_all"
	ProcessTRLocalChainCommitmentSubmit  = "local_commit_submit"
	ProcessTRLocalChainCommitmentPrepare = "local_commit_prepare"
	ProcessTRLocalCHainCommitmentRead    = "local_commit_read"
	ProcessTRLocalChainCommitmentReadAll = "local_commit_read_all"
	ProcessTROrgChainSubmit              = "org_submit"
	ProcessTROrgChainPrepare             = "org_prepare"
	ProcessTROrgChainRead                = "org_read"
	ProcessTROrgChainReadAll             = "org_read_all"

	ProcessCEMerkleProofVerify                       = "merkle_proof_verify"
	ProcessCEMerkleProofMerge                        = "merkle_proof_merge"
	ProcessCESummarizeMerkleProofVerificationResults = "summarize_merkle_proof"
	ProcessCEVerifyCommitments                       = "verify_commit"
)

var PhaseProcessMap = map[string][]string{
	PhaseInitialization: {
		ProcessINDefault,
		ProcessINRandGen,
	},
	PhaseTransactionRecord: {
		ProcessTRCommitment,
		ProcessTRMerkleProofGen,
		ProcessTRLocalChainSubmit,
		ProcessTRLocalChainPrepare,
		ProcessTRLocalChainRead,
		ProcessTRLocalChainReadAll,
		ProcessTRLocalChainCommitmentSubmit,
		ProcessTRLocalChainCommitmentPrepare,
		ProcessTRLocalCHainCommitmentRead,
		ProcessTRLocalChainCommitmentReadAll,
		ProcessTROrgChainSubmit,
		ProcessTROrgChainPrepare,
		ProcessTROrgChainRead,
		ProcessTROrgChainReadAll,
	},
	PhaseConsistencyExamination: {
		ProcessCEMerkleProofVerify,
		ProcessCEMerkleProofMerge,
		ProcessCESummarizeMerkleProofVerificationResults,
		ProcessCEVerifyCommitments,
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

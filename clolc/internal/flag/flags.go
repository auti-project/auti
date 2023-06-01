package flag

const (
	PhaseInitialization         = "in"
	PhaseTransactionRecord      = "tr"
	PhaseConsistencyExamination = "ce"
	PhaseResultVerification     = "rv"

	ProcessTRLocalChainSubmit     = "local_submit"
	ProcessTRLocalChainPrepare    = "local_prepare"
	ProcessTRLocalChainRead       = "local_read"
	ProcessTRLocalChainReadAll    = "local_read_all"
	ProcessTRCommitment           = "commit"
	ProcessTRAccumulate           = "accumulate"
	ProcessTROrgChainSubmit       = "org_submit"
	ProcessTROrgChainPrepare      = "org_prepare"
	ProcessTROrgChainRead         = "org_read"
	ProcessTROrgChainReadAll      = "org_read_all"
	ProcessCEAccumulateCommitment = "acc_commit"
	ProcessCEComputeB             = "cal_b"
	ProcessCEComputeC             = "cal_c"
	ProcessCEComputeD             = "cal_d"
	ProcessCEEncrypt              = "encrypt"
	ProcessCEAudChainSubmit       = "aud_submit"
	ProcessCEAudChainPrepare      = "aud_prepare"
	ProcessCEAudChainRead         = "aud_read"
	ProcessCEAudChainReadAll      = "aud_read_all"
	ProcessCECheck                = "check"
)

var (
	PhaseProcessMap = map[string][]string{
		PhaseInitialization: {},
		PhaseTransactionRecord: {
			ProcessTRLocalChainSubmit,
			ProcessTRLocalChainPrepare,
			ProcessTRLocalChainRead,
			ProcessTRLocalChainReadAll,
			ProcessTRCommitment,
			ProcessTRAccumulate,
			ProcessTROrgChainSubmit,
			ProcessTROrgChainPrepare,
			ProcessTROrgChainRead,
			ProcessTROrgChainReadAll,
		},
		PhaseConsistencyExamination: {
			ProcessCEAccumulateCommitment,
			ProcessCEComputeB,
			ProcessCEComputeC,
			ProcessCEComputeD,
			ProcessCEEncrypt,
			ProcessCEAudChainSubmit,
			ProcessCEAudChainPrepare,
			ProcessCEAudChainRead,
			ProcessCEAudChainReadAll,
			ProcessCECheck,
		},
		PhaseResultVerification: {},
	}
)

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

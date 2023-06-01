package constants

const (
	FlagPhaseInitialization         = "in"
	FlagPhaseTransactionRecord      = "tr"
	FlagPhaseConsistencyExamination = "ce"
	FlagPhaseResultVerification     = "rv"

	FlagProcessLocalChainSubmit     = "local_submit"
	FlagProcessLocalChainPrepare    = "local_prepare"
	FlagProcessLocalChainRead       = "local_read"
	FlagProcessLocalChainReadAll    = "local_read_all"
	FlagProcessCommitment           = "commit"
	FlagProcessAccumulate           = "accumulate"
	FlagProcessOrgChainPrepare      = "org_prepare"
	FlagProcessOrgChainSubmit       = "org_submit"
	FlagProcessOrgChainRead         = "org_read"
	FlagProcessOrgChainReadAll      = "org_read_all"
	FlagProcessAccumulateCommitment = "acc_commit"
	FlagProcessComputeB             = "cal_b"
	FlagProcessComputeC             = "cal_c"
	FlagProcessComputeD             = "cal_d"
	FlagProcessEncrypt              = "encrypt"
)

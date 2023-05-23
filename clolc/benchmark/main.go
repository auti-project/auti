package main

import (
	"flag"
)

func main() {
	benchPhasePtr := flag.String("phase", "i", "Benchmark to run")
	numOrgPtr := flag.Int("numOrg", 2, "Number of organizations")
	numIterPtr := flag.Int("numIter", 10, "Number of iterations")
	flag.Parse()
	switch *benchPhasePtr {
	case "i":
		benchCLOLCInitializeEpoch(*numOrgPtr, *numIterPtr)
	}
}

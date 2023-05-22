package main

import (
	"flag"
	"fmt"
)

func main() {
	benchPhasePtr := flag.String("phase", "i", "Benchmark to run")
	numOrgPtr := flag.Int("numOrg", 2, "Number of organizations")
	numIterPtr := flag.Int("numIter", 10, "Number of iterations")
	flag.Parse()
	fmt.Println("Benchmarking", *benchPhasePtr)
	switch *benchPhasePtr {
	case "i":
		benchCLOLCInitializeEpoch(*numOrgPtr, *numIterPtr)
	}
}

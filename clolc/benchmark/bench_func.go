package main

import (
	"fmt"
	"time"

	"github.com/auti-project/auti/clolc/benchmark/localchain"
	"github.com/auti-project/auti/internal/auditor"
	"github.com/auti-project/auti/internal/committee"
	"github.com/auti-project/auti/internal/organization"
)

func generateEntities(numOrganizations int) (*committee.Committee, []*auditor.Auditor, []*organization.Organization) {
	organizations := make([]*organization.Organization, numOrganizations)
	for i := 0; i < numOrganizations; i++ {
		organizations[i] = organization.New("org" + string(rune(i)))
	}
	auditors := make([]*auditor.Auditor, numOrganizations)
	for i := 0; i < numOrganizations; i++ {
		auditors[i] = auditor.New("aud"+string(rune(i)), []*organization.Organization{organizations[i]})
	}
	com := committee.New("com", auditors)
	return com, auditors, organizations
}

func benchInitializeEpoch(numOrganizations, iterations int) {
	fmt.Println("CLOLC initialize epoch")
	fmt.Println("Number of organizations:", numOrganizations)
	fmt.Println("Number of iterations:", iterations)
	fmt.Println("Elapsed times (ms):")
	runningTimes := make([]time.Duration, iterations)
	avgTime := time.Duration(0)
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(numOrganizations)
		startTime := time.Now()
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			panic(err)
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Println(elapsed.Milliseconds())
		runningTimes[i] = elapsed
		avgTime += elapsed
	}
	fmt.Printf("Average time (ms): %d\n\n", avgTime.Milliseconds()/int64(iterations))
}

func benchTransactionRecordSubmitTX(numTXs, iterations int) {
	fmt.Println("CLOLC transaction record submit transaction")
	fmt.Println("Number of transactions:", numTXs)
	fmt.Println("Number of iterations:", iterations)
	fmt.Println("Elapsed times (ms):")
	runningTimes := make([]time.Duration, iterations)
	avgTime := time.Duration(0)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		txIDs, err := localchain.BenchSubmitTX(numTXs)
		if err != nil {
			panic(err)
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Println(elapsed.Milliseconds())
		runningTimes[i] = elapsed
		avgTime += elapsed
		if err = localchain.SaveTXIDs(txIDs); err != nil {
			panic(err)
		}
	}
	fmt.Printf("Average time (ms): %d\n\n", avgTime.Milliseconds()/int64(iterations))
}

func benchTransactionRecordReadTX(iterations int) {
	fmt.Println("CLOLC transaction record read transaction")
	fmt.Println("Number of iterations:", iterations)
	fmt.Println("Elapsed times (ms):")
	runningTimes := make([]time.Duration, iterations)
	avgTime := time.Duration(0)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		if err := localchain.BenchReadTX(); err != nil {
			panic(err)
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Println(elapsed.Milliseconds())
		runningTimes[i] = elapsed
		avgTime += elapsed
	}
	fmt.Printf("Average time (ms): %d\n\n", avgTime.Milliseconds()/int64(iterations))
}

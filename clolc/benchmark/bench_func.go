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
	}
}

func benchTransactionRecordSubmitTX(numTXs, iterations int) {
	fmt.Println("CLOLC transaction record submit transaction")
	fmt.Println("Number of transactions:", numTXs)
	fmt.Println("Number of iterations:", iterations)
	fmt.Println("Elapsed times (ms):")
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		_, err := localchain.BenchSubmitTX(numTXs)
		if err != nil {
			panic(err)
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Println(elapsed.Milliseconds())
	}
	fmt.Println()
}

func benchTransactionRecordReadTX(numTotalTXs, iterations int) {
	fmt.Println("CLOLC transaction record read transaction")
	fmt.Println("Number of transactions:", numTotalTXs)
	fmt.Println("Number of iterations:", iterations)
	fmt.Println("Elapsed times (ms):")
	txIDs, err := localchain.BenchSubmitTX(numTotalTXs)
	if err != nil {
		panic(err)
	}
	if err = localchain.SaveTXIDs(txIDs); err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		if err := localchain.BenchReadTX(); err != nil {
			panic(err)
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Println(elapsed.Milliseconds())
	}
	fmt.Println()
}

func benchTransactionRecordReadAllTXs(numTotalTXs, iterations int) {
	fmt.Println("CLOLC transaction record read all transactions")
	fmt.Println("Number of transactions:", numTotalTXs)
	fmt.Println("Number of iterations:", iterations)
	fmt.Println("Elapsed times (ms):")
	txIDs, err := localchain.BenchSubmitTX(numTotalTXs)
	if err != nil {
		panic(err)
	}
	if err = localchain.SaveTXIDs(txIDs); err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		if err := localchain.BenchReadAllTXs(); err != nil {
			panic(err)
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Println(elapsed.Milliseconds())
	}
	fmt.Println()
}

package main

import (
	"fmt"
	"time"

	"github.com/auti-project/auti/core/auditor"
	"github.com/auti-project/auti/core/committee"
	"github.com/auti-project/auti/core/organization"
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
	runningTimes := make([]time.Duration, iterations)
	avgTime := time.Duration(0)
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
		runningTimes[i] = elapsed
		avgTime += elapsed
	}
	fmt.Println("Average time (ms):", avgTime.Milliseconds()/int64(iterations))
}

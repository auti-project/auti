package benchfeature

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/auti-project/auti/benchmark/timecounter"
	closcaud "github.com/auti-project/auti/internal/auditor/closc"
	closccom "github.com/auti-project/auti/internal/committee/closc"
	"github.com/auti-project/auti/internal/constants"
	closcorg "github.com/auti-project/auti/internal/organization/closc"
)

func generateEntities(numOrganizations int) (*closccom.Committee, []*closcaud.Auditor, []*closcorg.Organization) {
	organizations := make([]*closcorg.Organization, numOrganizations)
	for i := 0; i < numOrganizations; i++ {
		organizations[i] = closcorg.New("org" + string(rune(i)))
	}
	auditors := make([]*closcaud.Auditor, numOrganizations)
	for i := 0; i < numOrganizations; i++ {
		auditors[i] = closcaud.New("aud"+string(rune(i)), []*closcorg.Organization{organizations[i]})
	}
	com := closccom.New("com", auditors)
	return com, auditors, organizations
}

func InitializeEpoch(numOrganizations, iterations int) error {
	fmt.Println("[CLOSC-IN] Default")
	fmt.Printf("Num Org: %d, Num iter: %d\n", numOrganizations, iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, _ := generateEntities(numOrganizations)
		startTime := time.Now()
		err := com.InitializeEpoch(auditors)
		if err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func InitializeRandGen(num, iterations int) error {
	fmt.Print("[CLOSC-IN] Random generation")
	fmt.Printf("Num: %d, Num iter: %d\n", num, iterations)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		randByteList := make([][]byte, num)
		for j := 0; j < num; j++ {
			randByteList[j] = make([]byte, constants.SecurityParameterBytes)
			_, err := rand.Read(randByteList[j])
			if err != nil {
				return err
			}
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

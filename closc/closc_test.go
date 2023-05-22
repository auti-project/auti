package main

import (
	"testing"
)

func BenchmarkCLOSCInitializeEpoch(b *testing.B) {
	com, auditors, organizations := generateEntities(10)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			b.Error(err)
		}
	}
}

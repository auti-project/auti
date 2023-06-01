package benchfeature

import (
	"fmt"
	"time"
)

func printTime(elapsed time.Duration) {
	if elapsed.Milliseconds() <= 1 {
		fmt.Printf("Elapsed time: %d ns\n", elapsed.Nanoseconds())
	} else {
		fmt.Printf("Elapsed time: %d ms\n", elapsed.Milliseconds())
	}
}

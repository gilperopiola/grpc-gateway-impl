package benchmorkar

import (
	"fmt"
	"sync"
	"time"
)

// Lmao I'm seeing this after like a year of writing it and now I wonder why
// didn't I just use the testing package.

func Kumpare(fnsToBenchmork []func(int) string, runs int) {
	fmt.Printf("\n\n ~ [BENCHMORKAR] ~ \n")

	var mutex sync.Mutex
	var elapsedPerFn = make([]int64, len(fnsToBenchmork))

	for run := 1; run <= runs; run++ {
		for i, fn := range fnsToBenchmork {
			start := time.Now()
			fn(run)
			elapsed := time.Since(start)

			mutex.Lock()
			elapsedPerFn[i] += elapsed.Nanoseconds()
			mutex.Unlock()
		}
	}

	fmt.Println()

	for i := range fnsToBenchmork {
		fmt.Printf("Function %d: %dms ~ %dns\n", i, elapsedPerFn[i]/1000000, elapsedPerFn[i])
	}
}

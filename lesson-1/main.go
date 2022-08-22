package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	numbers := []int{}
	for i := 0; i < 1000; i++ {
		numbers = append(numbers, i)
	}

	t := time.Now()
	sum := add(numbers)
	fmt.Printf("Sequential Add, Sum: %d, Time Taken: %s\n", sum, time.Since(t))

	t = time.Now()
	sum = addConcurrent(numbers)
	fmt.Printf("Concurrent Add, Sum: %d, Time Taken: %s\n", sum, time.Since(t))
}

func add(numbers []int) int64 {
	var sum int64
	for _, n := range numbers {
		sum += int64(n)
	}
	return sum
}

func addConcurrent(numbers []int) int64 {
	numOfCores := runtime.NumCPU()
	runtime.GOMAXPROCS(numOfCores)

	var sum int64
	max := len(numbers)

	sizeOfParts := max / numOfCores

	var wg sync.WaitGroup

	for i := 0; i < numOfCores; i++ {
		start := i * sizeOfParts
		end := start + sizeOfParts
		part := numbers[start:end]

		wg.Add(1)
		go func(nums []int) {
			defer wg.Done()

			var partSum int64

			for _, n := range nums {
				partSum += int64(n)
			}

			atomic.AddInt64(&sum, partSum)
		}(part)
	}

	wg.Wait()
	return sum
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"
)

const DataFile = "loremipsum.txt"

// Return the word frequencies of the text argument.
func WordCount(text string) map[string]int {

	freqs := make(map[string]int)
	s := regexp.MustCompile("[\\s\\.\\,]").Split(strings.ToLower(text), -1)
	for _, f := range s {
		_, yes := freqs[f]
		if yes {
			freqs[f] += 1
		} else {
			freqs[f] = 1
		}
	}
	return freqs
}

// Benchmark how long it takes to count word frequencies in text numRuns times.
//
// Return the total time elapsed.
func benchmark(text string, numRuns int) int64 {
	start := time.Now()
	for i := 0; i < numRuns; i++ {
		WordCount(text)
	}
	runtimeMillis := time.Since(start).Nanoseconds() / 1e6

	return runtimeMillis
}

// Print the results of a benchmark
func printResults(runtimeMillis int64, numRuns int) {
	fmt.Printf("amount of runs: %d\n", numRuns)
	fmt.Printf("total time: %d ms\n", runtimeMillis)
	average := float64(runtimeMillis) / float64(numRuns)
	fmt.Printf("average time/run: %.2f ms\n", average)
}

func main() {
	// read in DataFile as a string called data
	a, b := ioutil.ReadFile(DataFile)
	if b != nil {
		log.Fatal(b)
	}

	d := string(a)
	fmt.Printf("%#v", WordCount(string(d)))

	numRuns := 100
	runtimeMillis := benchmark(string(d), numRuns)
	printResults(runtimeMillis, numRuns)
}

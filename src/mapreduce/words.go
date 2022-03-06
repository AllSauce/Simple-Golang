package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"
)

const DataFile = "loremipsum.txt"

// Return the word frequencies of the text argument.
func WordCount(text string) map[string]int {
	//Lower case letters
	freqs := make(map[string]int)
	ch := make(chan map[string]int)
	wg := new(sync.WaitGroup)
	wg1 := new(sync.WaitGroup)
	s := regexp.MustCompile("[\\.\\,\\s]").Split(strings.ToLower(text), -1)
	wg.Add(6)
	wg1.Add(1)

	//Temporary map for each of 6 goroutines that then get sent to channel ch
	for i := 0; i < 6; i++ {
		go func(i int) {
			m := make(map[string]int)
			for _, f := range s[(len(s)/6)*(i) : (len(s)/6)*(i+1)] {
				if f != "" {
					m[f]++
				}
			}
			ch <- m
			wg.Done()
		}(i)
	}

	//Add one to value for each same word in channel
	go func() {
		for m := range ch {
			for s, k := range m {
				if s != "" {
					freqs[s] += k
				}
			}
		}
		wg1.Done()
	}()

	wg.Wait()
	close(ch)
	wg1.Wait()
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

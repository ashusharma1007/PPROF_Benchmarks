package goroutineleak

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
)

func GoroutineCount() int {
	return runtime.NumGoroutine()
}

// LeakyFileReader spawns one goroutine per line but reads only 3 results.
// The channel is unbuffered, so every other goroutine blocks forever on its
// send and is never collected. /debug/pprof/goroutine shows them all parked
// on the same line.
func LeakyFileReader(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	results := make(chan string)

	for _, line := range lines {
		go func(l string) {
			results <- strings.ToUpper(l)
		}(line)
	}

	var output []string
	for i := 0; i < 3 && i < len(lines); i++ {
		output = append(output, <-results)
	}

	return output, nil
}

// SafeFileReader fixes the leak with three changes: a buffered channel so no
// send blocks, a WaitGroup to know when the senders are done, and a close +
// range so the receiver drains every result and terminates.
func SafeFileReader(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	results := make(chan string, len(lines))
	var wg sync.WaitGroup

	for _, line := range lines {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			results <- strings.ToUpper(l)
		}(line)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var output []string
	for result := range results {
		output = append(output, result)
	}

	return output, nil
}

// DemoLeak prints the goroutine count before and after a leaky read.
func DemoLeak() {
	tmpFile := createTempFile(50)
	defer os.Remove(tmpFile)

	before := GoroutineCount()
	fmt.Printf("Goroutines BEFORE leaky read: %d\n", before)

	results, _ := LeakyFileReader(tmpFile)
	fmt.Printf("Got %d results from leaky reader\n", len(results))

	after := GoroutineCount()
	fmt.Printf("Goroutines AFTER leaky read:  %d\n", after)
	fmt.Printf("LEAKED goroutines:            %d\n\n", after-before)
}

// DemoSafe prints the same counts for the non-leaking version. The "before"
// count is already inflated by whatever DemoLeak leaked: leaks outlive the
// function that created them.
func DemoSafe() {
	tmpFile := createTempFile(50)
	defer os.Remove(tmpFile)

	before := GoroutineCount()
	fmt.Printf("Goroutines BEFORE safe read: %d\n", before)

	results, _ := SafeFileReader(tmpFile)
	fmt.Printf("Got %d results from safe reader\n", len(results))

	after := GoroutineCount()
	fmt.Printf("Goroutines AFTER safe read:  %d\n", after)
	fmt.Printf("LEAKED goroutines:           %d\n\n", after-before)
}

func createTempFile(n int) string {
	f, err := os.CreateTemp("", "goroutine-leak-demo-*.txt")
	if err != nil {
		panic(err)
	}
	for i := 0; i < n; i++ {
		fmt.Fprintf(f, "line number %d with some data to process\n", i)
	}
	f.Close()
	return f.Name()
}

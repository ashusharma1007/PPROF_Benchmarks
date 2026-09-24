# Go pprof & Benchmarking Examples

Hands-on examples to learn Go profiling (`pprof`) and benchmarking (`testing.B`) using simple algorithms and data structures.

## What's Inside

| Package | What It Teaches | Key Insight |
|---------|----------------|-------------|
| `sorting/` | CPU profiling | BubbleSort O(n²) vs QuickSort O(n log n) — pprof shows which function eats CPU |
| `datastructures/` | Memory profiling | Slice vs Map vs LinkedList — different allocation patterns |
| `goroutineleak/` | Goroutine profiling | Leaky vs safe goroutine patterns — detect with `runtime.NumGoroutine()` and pprof |
| `memalloc/` | Heap analysis | String concat, slice growth, heap escape — where your memory goes |

## Prerequisites

- Go 1.21+ installed
- `graphviz` for pprof visualization: `sudo dnf install graphviz` (Fedora) or `sudo apt install graphviz` (Ubuntu)

## Quick Start

### Step 1: Run the Benchmarks

```bash
cd pprof-benchmark-examples

# Run ALL benchmarks with memory stats
go test -bench=. -benchmem ./...

# Run a specific package
go test -bench=. -benchmem ./sorting/
go test -bench=. -benchmem ./memalloc/
```

### Step 2: Read the Benchmark Output

```
BenchmarkBubbleSort-16    1947      520477 ns/op    8192 B/op    1 allocs/op
BenchmarkQuickSort-16    77016       13973 ns/op    8192 B/op    1 allocs/op
```

| Column | Meaning |
|--------|---------|
| `-16` | Number of CPU cores used |
| `1947` | Number of iterations the benchmark ran |
| `520477 ns/op` | Nanoseconds per operation (lower = faster) |
| `8192 B/op` | Bytes allocated per operation |
| `1 allocs/op` | Heap allocations per operation |

### Step 3: Generate Profile Files

```bash
# Generate CPU and memory profiles from benchmarks
go test -bench=BenchmarkBubbleSort -cpuprofile=cpu.prof -memprofile=mem.prof ./sorting/

# Analyze CPU profile
go tool pprof cpu.prof

# Analyze memory profile
go tool pprof mem.prof
```

### Step 4: Use pprof Interactive Commands

Once inside the pprof interactive shell:

```
(pprof) top          # Show top functions by resource usage
(pprof) top10        # Top 10 functions
(pprof) list Bubble  # Show line-by-line cost for BubbleSort
(pprof) web          # Open flame graph in browser (needs graphviz)
(pprof) png          # Save flame graph as PNG
(pprof) quit         # Exit
```

### Step 5: HTTP pprof (Live Profiling)

```bash
# Start the demo server with pprof endpoint
go run main.go

# In another terminal — profile goroutines (find leaks!)
go tool pprof http://localhost:6060/debug/pprof/goroutine

# Profile heap memory
go tool pprof http://localhost:6060/debug/pprof/heap

# CPU profile for 30 seconds
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Or just open in browser
# http://localhost:6060/debug/pprof/
```

## Example Walkthroughs

### Goroutine Leak Detection

This is the most practical example. Run the tests to see the leak:

```bash
go test -v -run TestLeaky ./goroutineleak/
```

Output:
```
=== RUN   TestLeakyFileReader
    leak_test.go:36: Goroutines before: 2, after: 49, leaked: 47
    leak_test.go:41: CONFIRMED: 47 goroutines leaked!
--- PASS: TestLeakyFileReader
```

Then see the safe version has zero leaks:
```bash
go test -v -run TestSafe ./goroutineleak/
```

Output:
```
=== RUN   TestSafeFileReader
    leak_test.go:59: Goroutines before: 2, after: 2, leaked: 0
    leak_test.go:60: Processed 50 lines without leaking
--- PASS: TestSafeFileReader
```

**What causes the leak?** The leaky version spawns 50 goroutines to process file lines, but only reads 3 results from an unbuffered channel. The other 47 goroutines block forever trying to send — they never exit.

**The fix:** Use a buffered channel (capacity = number of goroutines) and a WaitGroup to drain all results. See `goroutineleak/leak.go` for both versions side by side.

### String Concatenation Performance

```bash
go test -bench=BenchmarkConcat -benchmem ./memalloc/
```

```
BenchmarkConcatWithPlus-16       159    9055150 ns/op   53164110 B/op   10001 allocs/op
BenchmarkConcatWithBuilder-16  59442      17619 ns/op      10240 B/op       1 allocs/op
BenchmarkConcatWithBytes-16   129614      11278 ns/op      10240 B/op       1 allocs/op
```

**Result:** `strings.Builder` is **500x faster** and uses **5000x less memory** than `+= ` concatenation. The `+` operator creates a new string copy every iteration.

### Sorting Algorithm Scaling

```bash
go test -bench=BenchmarkBubbleSortSizes -benchmem ./sorting/
go test -bench=BenchmarkQuickSortSizes -benchmem ./sorting/
```

Watch BubbleSort times grow quadratically (100x slower at 5000 vs 500 elements) while QuickSort grows linearly.

## Useful Benchmark Flags

```bash
# Run benchmarks for 5 seconds each (more stable results)
go test -bench=. -benchtime=5s ./sorting/

# Run 3 times for statistical comparison
go test -bench=. -count=3 ./sorting/

# Run only benchmarks matching a regex
go test -bench=BenchmarkQuickSort ./sorting/

# See compiler escape analysis (what goes to heap)
go build -gcflags="-m" ./memalloc/
```

## pprof Profile Types

| Profile | URL Path | What It Shows |
|---------|----------|---------------|
| CPU | `/debug/pprof/profile?seconds=30` | Where CPU time is spent |
| Heap | `/debug/pprof/heap` | Current memory usage |
| Allocs | `/debug/pprof/allocs` | All past allocations (even freed) |
| Goroutine | `/debug/pprof/goroutine` | All goroutine stacks (find leaks!) |
| Block | `/debug/pprof/block` | Where goroutines block (channels, mutexes) |
| Mutex | `/debug/pprof/mutex` | Mutex contention |

## Project Structure

```
pprof-benchmark-examples/
├── main.go                  # HTTP pprof server + demo runner
├── go.mod
├── sorting/
│   ├── sort.go              # BubbleSort vs QuickSort
│   └── sort_test.go         # Benchmarks with sub-benchmarks
├── datastructures/
│   ├── ds.go                # Slice vs Map vs LinkedList
│   └── ds_test.go           # Build + search benchmarks
├── goroutineleak/
│   ├── leak.go              # Leaky vs safe goroutine patterns
│   └── leak_test.go         # Tests that PROVE leaks + benchmarks
├── memalloc/
│   ├── alloc.go             # String concat, slice growth, heap escape
│   └── alloc_test.go        # Memory allocation benchmarks
├── CLAUDE.md                # Step-by-step learning guide
└── README.md                # This file
```
# PPROF_Benchmarks

package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"pprof-benchmark-examples/goroutineleak"
)

func main() {
	// The blank net/http/pprof import registers the profiling handlers on the
	// default mux. Bind to localhost only: these endpoints expose stack traces
	// and let any caller impose CPU load.
	go func() {
		log.Fatal(http.ListenAndServe("localhost:6060", nil))
	}()

	// Leaks 47 goroutines, then shows the version that doesn't.
	goroutineleak.DemoLeak()
	goroutineleak.DemoSafe()

	fmt.Printf("goroutines alive: %d (the leaked ones are still parked)\n\n", runtime.NumGoroutine())
	fmt.Println("profile from another terminal:")
	fmt.Println("  go tool pprof http://localhost:6060/debug/pprof/goroutine")
	fmt.Println("  go tool pprof -http=:8080 http://localhost:6060/debug/pprof/goroutine")
	fmt.Println("\nctrl-c to stop.")

	select {}
}

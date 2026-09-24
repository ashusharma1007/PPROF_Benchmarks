package goroutineleak

import (
	"os"
	"runtime"
	"testing"
	"time"
)

func TestLeakyFileReader(t *testing.T) {
	tmpFile := createTempFile(50)
	defer os.Remove(tmpFile)

	before := runtime.NumGoroutine()
	_, err := LeakyFileReader(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	after := runtime.NumGoroutine()

	leaked := after - before
	t.Logf("Goroutines before: %d, after: %d, leaked: %d", before, after, leaked)
	if leaked <= 0 {
		t.Error("Expected goroutine leak but none detected")
	}
	t.Logf("CONFIRMED: %d goroutines leaked!", leaked)
}

func TestSafeFileReader(t *testing.T) {
	tmpFile := createTempFile(50)
	defer os.Remove(tmpFile)

	before := runtime.NumGoroutine()
	results, err := SafeFileReader(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	after := runtime.NumGoroutine()

	leaked := after - before
	t.Logf("Goroutines before: %d, after: %d, leaked: %d", before, after, leaked)
	t.Logf("Processed %d lines without leaking", len(results))
	if leaked > 0 {
		t.Errorf("Expected no leak but %d goroutines leaked", leaked)
	}
}

func BenchmarkLeakyFileReader(b *testing.B) {
	tmpFile := createTempFile(100)
	defer os.Remove(tmpFile)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LeakyFileReader(tmpFile)
	}
}

func BenchmarkSafeFileReader(b *testing.B) {
	tmpFile := createTempFile(100)
	defer os.Remove(tmpFile)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SafeFileReader(tmpFile)
	}
}

package internal

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"
)

// SetupMemoryProfiling sets up memory profiling and returns a cleanup function
// that writes the heap profile when called.
// Note: We return a function instead of writing the profile directly because
// memory profiling captures a snapshot of the heap at the moment it's called.
// By deferring the returned function in main, we capture memory state at program
// exit (after all work is done), not at program start when nothing has happened yet.
func SetupMemoryProfiling(profilesDir string, timestamp time.Time) func() {
	memProfileFile := filepath.Join(profilesDir, fmt.Sprintf("mem_profile_%s.prof", timestamp.Format("2006-01-02_15-04-05")))

	return func() {
		mf, err := os.Create(memProfileFile)
		if err != nil {
			slog.Error("could not create memory profile", slog.Any("error", err))
			return
		}
		defer mf.Close()

		if err := pprof.WriteHeapProfile(mf); err != nil {
			slog.Error("could not write memory profile", slog.Any("error", err))
			return
		}
		slog.Info("memory profile written", slog.String("file", memProfileFile))
	}
}

package internal

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"
)

// StartCPUProfiling starts CPU profiling and returns a cleanup function
func StartCPUProfiling(profilesDir string, timestamp time.Time) (cleanup func(), err error) {
	// Create profiles directory if it doesn't exist
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create profiles directory: %w", err)
	}

	cpuProfileFile := filepath.Join(profilesDir, fmt.Sprintf("cpu_profile_%s.prof", timestamp.Format("2006-01-02_15-04-05")))
	f, err := os.Create(cpuProfileFile)
	if err != nil {
		return nil, fmt.Errorf("could not create CPU profile: %w", err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		return nil, fmt.Errorf("could not start CPU profile: %w", err)
	}

	slog.Info("CPU profiling enabled", slog.String("file", cpuProfileFile))

	cleanup = func() {
		pprof.StopCPUProfile()
		f.Close()
	}

	return cleanup, nil
}

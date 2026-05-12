package cli

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

type CLIService struct{}

type CLIOptions struct {
	FilePath string
	Format   string
	Timeout  time.Duration
}

func NewCLIService() CLIService {
	return CLIService{}
}

func (cs *CLIService) GetCLIOptions() (*CLIOptions, error) {
	filePath := flag.String("file", "", "path to the file with target URLs")
	format := flag.String("format", "text", "output format: text or json")
	timeout := flag.String("timeout", "30s", "execution timeout (e.g., 30s, 5m)")
	flag.Parse()

	if *filePath == "" {
		return nil, fmt.Errorf("error: -file flag is required")
	}

	requestedFormat := strings.ToLower(strings.TrimSpace(*format))
	supportedFormats := map[string]bool{
		"text": true,
		"json": true,
	}

	if !supportedFormats[requestedFormat] {
		return nil, fmt.Errorf("error: unsupported output format: %s", *format)
	}

	// Parse timeout duration
	duration, err := time.ParseDuration(*timeout)
	if err != nil {
		return nil, fmt.Errorf("error: invalid timeout value: %s (expected format like 30s, 5m, 1h)", *timeout)
	}

	if duration <= 0 {
		return nil, fmt.Errorf("error: timeout must be greater than zero")
	}

	return &CLIOptions{
		FilePath: *filePath,
		Format:   requestedFormat,
		Timeout:  duration,
	}, nil
}

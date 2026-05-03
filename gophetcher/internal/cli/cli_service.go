package cli

import (
	"flag"
	"fmt"
	"strings"
)

type CLIService struct{}

type CLIOptions struct {
	FilePath string
	Format   string
}

func NewCLIService() CLIService {
	return CLIService{}
}

func (cs *CLIService) GetCLIOptions() (*CLIOptions, error) {
	filePath := flag.String("file", "", "path to the file with target URLs")
	format := flag.String("format", "text", "output format: text or json")
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

	return &CLIOptions{
		FilePath: *filePath,
		Format:   requestedFormat,
	}, nil
}

package internal

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mdobak/go-xerrors"
)

// GetTxtFiles returns a list of all .txt file paths in the given directory
func GetTxtFiles(directoryPath string) ([]string, error) {
	var txtFiles []string

	// Check if directory exists
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		return nil, xerrors.Newf("directory does not exist: %s: %w", directoryPath, err)
	}

	// Walk through the directory
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if it's a file and has .txt extension
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".txt" {
			txtFiles = append(txtFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, xerrors.Newf("error walking directory: %w", err)
	}

	return txtFiles, nil
}

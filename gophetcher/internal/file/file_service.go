package file

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

func (fs *FileService) GetTargetFilePath() (*string, error) {
	filePath := flag.String("file", "", "path to the file with target URLs")
	flag.Parse()

	if *filePath == "" {
		return nil, fmt.Errorf("error: -file flag is required")
	}

	return filePath, nil
}

func (fs *FileService) ValidateFile(filePath string) error {
	if filepath.Ext(filePath) != ".txt" {
		return fmt.Errorf("error: only .txt files are allowed")
	}

	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("error: file does not exist: %s", filePath)
	}
	if err != nil {
		return fmt.Errorf("error: cannot access file: %v", err)
	}
	if info.IsDir() {
		return fmt.Errorf("error: path is a directory, not a file: %s", filePath)
	}

	return nil
}

func (fs *FileService) ReadURLsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error: cannot open file: %v", err)
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error: cannot read file: %v", err)
	}

	return urls, nil
}

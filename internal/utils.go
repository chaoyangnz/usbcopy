package internal

import (
	"fmt"
	"io"
	logging "log"
	"os"
	"path/filepath"
	"strings"
)

func log(format string, args ...any) {
	logging.Printf(format+"\n", args...)
}

func filter(extension string, whitelist []string) bool {
	for _, ext := range whitelist {
		if ext == extension {
			return true
		}
	}
	return false
}

func interpolate(s string, vars []string, vals []string) string {
	str := s
	for i, _ := range vars {
		str = strings.ReplaceAll(str, vars[i], vals[i])
	}
	return filepath.Clean(str)
}

func extractExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ""
	}
	return strings.ToUpper(ext[1:])
}

func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Couldn't open dest file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("Couldn't copy to dest from source: %v", err)
	}

	inputFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801

	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't remove source file: %v", err)
	}
	return nil
}

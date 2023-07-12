package utils

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/leonardobiffi/kse/package/entities"
	"gopkg.in/yaml.v2"
)

func ReadStdin(rd io.Reader) []byte {
	var output []byte
	reader := bufio.NewReader(rd)
	for {
		input, err := reader.ReadByte()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}
	return output
}

func ReadFile(file string) ([]byte, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// Update file content
func UpdateFile(file string, content []byte) error {
	os.Remove(file)

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		return err
	}

	return nil
}

// find secret files in a directory
func FindFiles(dir string) ([]string, error) {
	var files []string

	if dir == "." {
		dir, _ = os.Getwd()
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// ignore dot files
		if info.Name()[0] == '.' {
			return filepath.SkipDir
		}

		if !info.IsDir() && IsSecret(path) {
			files = append(files, sanitizeDir(path))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// IsSecret check if file is a kubernetes secret
func IsSecret(file string) bool {
	content, err := os.ReadFile(file)
	if err != nil {
		return false
	}

	var secret entities.Secret
	err = yaml.Unmarshal(content, &secret)
	if err != nil {
		return false
	}

	if secret.Kind == "Secret" {
		return true
	}

	return false
}

// remove current directory from path
func sanitizeDir(dir string) string {
	curDir, _ := os.Getwd()
	return strings.Replace(dir, curDir+"/", "", 1)
}

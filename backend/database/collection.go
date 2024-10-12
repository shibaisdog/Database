package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ParseJSON[T any](f string, t *T) error {
	file, err := ReadFile(f)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(file))
	return decoder.Decode(t)
}

func ByteJSON[T any](t *T) ([]byte, error) {
	jsonData, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func ReadFile(f string) ([]byte, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil
}

func WriteFile(f string, data []byte) error {
	file, err := os.Create(f)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func DeleteFile(f string) error {
	err := os.Remove(f)
	if err != nil {
		return err
	}
	return nil
}

func ListFiles(f string) ([]string, error) {
	var files []string
	items, err := os.ReadDir(f)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if !item.IsDir() {
			fullPath := filepath.Join(f, item.Name())
			files = append(files, fullPath)
		}
	}
	return files, nil
}

func FileName(f ...string) string {
	name := ""
	for i, n := range f {
		if i != len(f)-1 {
			name += n + "/"
		} else {
			name += n
		}
	}
	return name
}

package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type DataJson struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

func readJSON(filename string) ([]DataJson, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var data []DataJson
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return data, nil
}

func writeJSON(filename string, data []DataJson) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

func CheckHash(storageFile string, hashString string) bool {
	data, err := readJSON(storageFile)
	if err != nil {
		return false
	}

	// Ищем хеш и фильтруем массив
	filteredData := make([]DataJson, 0)
	found := false

	for _, item := range data {
		if item.Hash == hashString {
			found = true
			continue // Пропускаем найденную запись
		}
		filteredData = append(filteredData, item)
	}

	if found {
		if err := writeJSON(storageFile, filteredData); err != nil {
			return false
		}
	}

	return found
}

func AddEmailHash(storageFile string, email string, hashString string) error {
	data, err := readJSON(storageFile)
	if err != nil {
		return err
	}
	data = append(data, DataJson{
		Email: email,
		Hash:  hashString,
	})
	writeJSON(storageFile, data)
	return nil
}

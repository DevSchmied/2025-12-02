package storage

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

// FileReader is an interface for abstracting file read/write operations (DI).
type FileReader interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
}

// OSFileReader is a FileReader implementation based on os.
type OSFileReader struct{}

// ReadFile reads a file using os.ReadFile.
func (oSFileReader OSFileReader) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile writes a file using os.WriteFile.
func (oSFileReader OSFileReader) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

// Storage stores data in memory and manages reading/writing a JSON file.
type Storage struct {
	mu          sync.Mutex
	Data        map[int]map[string]string
	LastLinkNum int
	filePath    string
	reader      FileReader
}

// storageFileData is a DTO for serializing data to JSON.
type storageFileData struct {
	LastLinkNum int                       `json:"last_link_num"`
	Data        map[int]map[string]string `json:"data"`
}

// NewStorage creates a storage object and immediately loads data from disk.
func NewStorage(filePath string, reader FileReader) (*Storage, error) {
	if reader == nil {
		reader = OSFileReader{}
	}

	strg := &Storage{
		Data:        make(map[int]map[string]string),
		LastLinkNum: 0,
		filePath:    filePath,
		reader:      reader,
	}

	if err := strg.LoadFromDisk(); err != nil {
		log.Printf("storage file error: %v — starting with empty storage", err)
	}

	return strg, nil
}

// LoadFromDisk loads storage data from a JSON file.
func (strg *Storage) LoadFromDisk() error {
	// Lock access to avoid data races
	strg.mu.Lock()
	defer strg.mu.Unlock()

	fileData, err := strg.reader.ReadFile(strg.filePath)
	if err != nil {
		// File does not exist — initialize empty state
		if errors.Is(err, os.ErrNotExist) {
			strg.Data = make(map[int]map[string]string)
			strg.LastLinkNum = 0
			return nil
		}
		return err
	}

	// Convert the JSON file into a Go structure to extract LastLinkNum and Data
	var parsed storageFileData
	if err := json.Unmarshal(fileData, &parsed); err != nil {
		return err
	}

	if parsed.Data == nil {
		parsed.Data = make(map[int]map[string]string)
	}

	strg.Data = parsed.Data
	strg.LastLinkNum = parsed.LastLinkNum

	return nil
}

// SaveToDisk saves the current storage state to a JSON file.
func (strg *Storage) SaveToDisk() error {
	strg.mu.Lock()
	defer strg.mu.Unlock()

	fileData := &storageFileData{
		LastLinkNum: strg.LastLinkNum,
		Data:        strg.Data,
	}

	// Convert to JSON with indentation
	encoded, err := json.MarshalIndent(fileData, "", "  ")
	if err != nil {
		return err
	}

	// Write to file (the file will be created if it does not exist)
	return strg.reader.WriteFile(strg.filePath, encoded, 0644)
}

// GenerateID increments the counter and returns a new number.
func (strg *Storage) GenerateID() int {
	strg.mu.Lock()
	defer strg.mu.Unlock()
	strg.LastLinkNum++
	return strg.LastLinkNum
}

// AddRecord stores a new result by ID.
func (strg *Storage) AddRecord(id int, data map[string]string) {
	strg.mu.Lock()
	defer strg.mu.Unlock()

	strg.Data[id] = data
}

// GetRecords returns data for the specified link group IDs.
func (strg *Storage) GetRecords(ids []int) map[int]map[string]string {
	strg.mu.Lock()
	defer strg.mu.Unlock()

	out := make(map[int]map[string]string)

	for _, id := range ids {
		if val, ok := strg.Data[id]; ok {
			out[id] = val
		}
	}

	return out
}

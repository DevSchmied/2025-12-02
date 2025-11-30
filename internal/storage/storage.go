package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// Storage хранит данные в памяти и управляет чтением/записью JSON-файла.
type Storage struct {
	mu          sync.Mutex
	Data        map[int]map[string]string
	LastLinkNum int
	filePath    string
}

// storageFileData — DTO для сериализации данных в JSON.
type storageFileData struct {
	LastLinkNum int                       `json:"last_link_num"`
	Data        map[int]map[string]string `json:"data"`
}

// NewStorage создаёт объект хранилища и сразу загружает данные с диска.
func NewStorage(filePath string) (*Storage, error) {
	strg := &Storage{
		Data:        make(map[int]map[string]string),
		LastLinkNum: 0,
		filePath:    filePath,
	}

	if err := strg.LoadFromDisk(); err != nil {
		return nil, err
	}

	return strg, nil
}

// LoadFromDisk загружает данные хранилища из JSON-файла.
func (strg *Storage) LoadFromDisk() error {
	// Блокируем доступ, чтобы избежать гонок данных
	strg.mu.Lock()
	defer strg.mu.Unlock()

	fileData, err := os.ReadFile(strg.filePath)
	if err != nil {
		// Файл отсутствует — инициализируем пустое состояние
		if errors.Is(err, os.ErrNotExist) {
			strg.Data = make(map[int]map[string]string)
			strg.LastLinkNum = 0
			return nil
		}
		return err
	}

	// Превращаем JSON-файл в структуру Go, чтобы можно было извлечь LastLinkNum и Data
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

// SaveToDisk сохраняет текущее состояние хранилища в JSON-файл.
func (strg *Storage) SaveToDisk() error {
	strg.mu.Lock()
	defer strg.mu.Unlock()

	fileData := &storageFileData{
		LastLinkNum: strg.LastLinkNum,
		Data:        strg.Data,
	}

	// Преобразуем в JSON с отступами
	encoded, err := json.MarshalIndent(fileData, "", "  ")
	if err != nil {
		return err
	}

	// Записываем в файл (файл будет создан, если его нет)
	return os.WriteFile(strg.filePath, encoded, 0644)
}

// GenerateID увеличивает счётчик и возвращает новый номер.
func (strg *Storage) GenerateID() int {
	strg.mu.Lock()
	defer strg.mu.Unlock()
	strg.LastLinkNum++
	return strg.LastLinkNum
}

// AddRecord сохраняет новый результат по ID.
func (strg *Storage) AddRecord(id int, data map[string]string) {
	strg.mu.Lock()
	defer strg.mu.Unlock()

	strg.Data[id] = data
}

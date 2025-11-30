package main

import (
	"2025/internal/server"
	"2025/internal/storage"
	"fmt"
	"log"
	"sync"
)

func main() {
	fmt.Println("Test")

	addr := "localhost:8080"
	server := server.NewServer(addr)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server failed to start: %v\n", err)
		}
	}()

	jsonPath := "internal\\storage\\storage.json"
	storage, err := storage.NewStorage(jsonPath)
	_ = err
	err2 := storage.SaveToDisk()
	fmt.Println("TEST:", err2)

	wg.Wait()
}

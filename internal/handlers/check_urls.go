package handlers

import (
	"2025/internal/service"
	"2025/internal/storage"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CheckURLs проверяет URL через worker pool, сохраняет результаты и возвращает links_num.
func CheckURLs(storage *storage.Storage, tasks chan service.Task) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Links []string `json:"links"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}

		id := storage.GenerateID()

		// канал для получения результатов от воркеров
		resultCh := make(chan service.Result)
		// map для будущего ответа
		results := make(map[string]string)

		// каждую ссылку в worker pool
		for _, link := range req.Links {
			tasks <- service.Task{
				URL: link,
				Res: resultCh,
			}
		}

		// результаты из канала
		for range req.Links {
			res := <-resultCh
			results[res.URL] = res.Status
		}

		close(resultCh)

		// сохраняем
		storage.AddRecord(id, results)
		err := storage.SaveToDisk()
		if err != nil {
			log.Printf("failed to write storage to JSON file: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"links":     results,
			"links_num": id,
		})
	}
}

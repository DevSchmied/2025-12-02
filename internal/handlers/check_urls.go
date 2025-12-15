package handlers

import (
	"2025/internal/service"
	"2025/internal/storage"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// CheckURLs checks URLs via a worker pool, stores the results, and returns links_num.
func CheckURLs(
	storage *storage.Storage,
	tasks chan service.Task,
	wg *sync.WaitGroup,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Links []string `json:"links"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}

		id := storage.GenerateID()

		// Channel for receiving results from workers
		resultCh := make(chan service.Result)
		// Map for the future response
		results := make(map[string]string)

		// Send each link to the worker pool
		for _, link := range req.Links {
			wg.Add(1)
			tasks <- service.Task{
				URL: link,
				Res: resultCh,
			}
		}

		// Collect results from the channel
		for range req.Links {
			res := <-resultCh
			results[res.URL] = res.Status
		}

		close(resultCh)

		// Save results
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

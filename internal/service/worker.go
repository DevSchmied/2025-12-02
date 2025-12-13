package service

import (
	"2025/internal/check"
	"log"
	"sync"
)

type Result struct {
	URL    string
	Status string
}

type Task struct {
	URL string
	Res chan Result
}

// StartWorkerPool starts N workers
// that process tasks from the tasks channel
func StartWorkerPool(n int, tasks chan Task, wg *sync.WaitGroup) {
	for i := 1; i <= n; i++ {
		go func(workerId int) {
			// Workers read tasks from the channel until it is closed
			for task := range tasks {
				wg.Add(1)
				func() {
					defer wg.Done()
					log.Printf("worker %d processing %s", workerId, task.URL)

					if check.CheckLink(task.URL) {
						task.Res <- Result{URL: task.URL, Status: "available"}
					} else {
						task.Res <- Result{URL: task.URL, Status: "not available"}
					}
				}()
			}
		}(i)
	}
}

package handlers

import (
	"2025/internal/pdf"
	"2025/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MakePDF — хендлер, который принимает список номеров ссылок,
// формирует по ним PDF и отправляет его пользователю.
func MakePDF(storage *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ListNums []int `json:"links_list"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}

		data := storage.GetRecords(req.ListNums)

		pdfBytes, err := pdf.GeneratePDF(data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "pdf generation error"})
			return
		}

		// Указываем, что возвращаем файл PDF
		c.Header("Content-Type", "application/pdf")

		// Говорим браузеру скачать файл под именем report.pdf
		c.Header("Content-Disposition", "attachment; filename=report.pdf")

		// Отправляем PDF как бинарный ответ
		c.Data(http.StatusOK, "application/pdf", pdfBytes)
	}
}

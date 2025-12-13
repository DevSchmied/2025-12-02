package handlers

import (
	"2025/internal/pdf"
	"2025/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MakePDF is a handler that accepts a list of link set numbers,
// generates a PDF based on them, and sends it to the user.
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

		// Specify that a PDF file is being returned
		c.Header("Content-Type", "application/pdf")

		// Instruct the browser to download the file as report.pdf
		c.Header("Content-Disposition", "attachment; filename=report.pdf")

		// Send the PDF as a binary response
		c.Data(http.StatusOK, "application/pdf", pdfBytes)
	}
}

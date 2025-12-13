package pdf

import (
	"fmt"
	"path/filepath"

	"github.com/signintech/gopdf"
)

// GeneratePDF generates a PDF file containing a list of link groups and their statuses.
func GeneratePDF(data map[int]map[string]string) ([]byte, error) {
	var pdf gopdf.GoPdf

	// Create a new PDF document
	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
	})
	// Add a single page
	pdf.AddPage()

	// Load the font
	if err := pdf.AddTTFFont("DejaVu", filepath.Join("internal", "pdf", "DejaVuSans.ttf")); err != nil {
		return nil, fmt.Errorf("cannot load font: %v", err)
	}
	if err := pdf.SetFont("DejaVu", "", 14); err != nil {
		return nil, fmt.Errorf("cannot set font: %v", err)
	}

	// Initial vertical coordinate
	x := 30.0
	y := 40.0

	for id, links := range data {
		// Group header
		pdf.SetXY(x, y)
		pdf.Cell(nil, fmt.Sprintf("Group %d:", id))
		y += 18

		// Lines with URLs and statuses
		for url, status := range links {
			pdf.SetXY(x+10, y)
			pdf.Cell(nil, fmt.Sprintf("%s - %s", url, status))
			y += 15
		}

		// Spacing between groups
		y += 15
	}
	return pdf.GetBytesPdf(), nil
}

package pdf

import (
	"fmt"
	"path/filepath"

	"github.com/signintech/gopdf"
)

// GeneratePDF формирует PDF-файл со списком групп ссылок и их статусами.
func GeneratePDF(data map[int]map[string]string) ([]byte, error) {
	var pdf gopdf.GoPdf

	// Создаём новый PDF-документ
	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
	})
	// Добавляем одну страницу
	pdf.AddPage()

	// Подключаем шрифт
	if err := pdf.AddTTFFont("DejaVu", filepath.Join("internal", "pdf", "DejaVuSans.ttf")); err != nil {
		return nil, fmt.Errorf("cannot load font: %v", err)
	}
	if err := pdf.SetFont("DejaVu", "", 14); err != nil {
		return nil, fmt.Errorf("cannot set font: %v", err)
	}

	// Начальная координата по вертикали
	x := 30.0
	y := 40.0

	for id, links := range data {
		// Заголовок группы
		pdf.SetXY(x, y)
		pdf.Cell(nil, fmt.Sprintf("Group %d:", id))
		y += 18

		// Строки с URL и статусами
		for url, status := range links {
			pdf.SetXY(x+10, y)
			pdf.Cell(nil, fmt.Sprintf("%s - %s", url, status))
			y += 15
		}

		// Отступ между группами
		y += 15
	}
	return pdf.GetBytesPdf(), nil
}

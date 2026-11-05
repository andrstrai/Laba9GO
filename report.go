package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// SalesReport — итоговый отчет по продажам.
// Именно эта структура пишется в report.json.
type SalesReport struct {
	GeneratedAt   string                  `json:"generated_at"`            // когда сформирован отчет
	ProductsCount int                     `json:"products_count"`          // сколько товаров загружено
	SalesCount    int                     `json:"sales_count"`             // сколько записей о продажах загружено
	TotalRevenue  float64                 `json:"total_revenue"`           // общая выручка
	TopProducts   []ProductSales          `json:"top_products"`            // топ самых продаваемых товаров
	Categories    map[string]CategoryStat `json:"categories"`              // статистика по категориям
	UnknownSales  []Sale                  `json:"unknown_sales,omitempty"` // продажи без товара в каталоге
}

// BuildSalesReport собирает итоговый отчет из уже загруженных данных.
// topN — сколько товаров включать в топ (<= 0 — все).
func BuildSalesReport(products []Product, sales []Sale, topN int) SalesReport {
	report := SalesReport{
		GeneratedAt:   time.Now().Format(time.RFC3339),
		ProductsCount: len(products),
		SalesCount:    len(sales),
		TotalRevenue:  TotalRevenue(products, sales),
		TopProducts:   TopProducts(products, sales, topN),
		Categories:    CategoryStats(products, sales),
		UnknownSales:  UnknownSales(products, sales),
	}
	return report
}

// SaveSalesReport сериализует отчет в JSON с отступами и записывает его в файл.
func SaveSalesReport(report SalesReport, filename string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("не удалось сериализовать отчёт в JSON: %w", err)
	}

	// Добавляем перевод строки в конец — так файл удобнее читать и diff'ить в git.
	data = append(data, '\n')

	if err := os.WriteFile(filename, data, 0o644); err != nil {
		return fmt.Errorf("не удалось записать файл %q: %w", filename, err)
	}
	return nil
}

// WriteReport — обертка: собрать отчет и сразу сохранить.
func WriteReport(products []Product, sales []Sale, topN int, filename string) error {
	if len(products) == 0 {
		return fmt.Errorf("список товаров пуст: сначала загрузите JSON с продуктами")
	}
	if len(sales) == 0 {
		return fmt.Errorf("список продаж пуст: сначала загрузите CSV с продажами")
	}

	report := BuildSalesReport(products, sales, topN)
	if err := SaveSalesReport(report, filename); err != nil {
		return err
	}

	fmt.Printf("Отчет сохранён в %q\n", filename)
	fmt.Printf("Товаров: %d, продаж: %d, выручка: %.2f\n",
		report.ProductsCount, report.SalesCount, report.TotalRevenue)
	if len(report.UnknownSales) > 0 {
		fmt.Printf("Внимание: %d продаж ссылаются на несуществующие товары (см. поле unknown_sales)\n",
			len(report.UnknownSales))
	}
	return nil
}

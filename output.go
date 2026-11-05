package main

import (
	"bufio"
	"fmt"
	"sort"
	"strings"
)

const tableLine = "------------------------------------------------------------------"

// PrintProducts печатает список товаров таблицей.
// Используется после FilterByCategory / FilterByPrice.
func PrintProducts(products []Product) {
	if len(products) == 0 {
		fmt.Println("Товары не найдены.")
		return
	}
	fmt.Printf("%-5s %-25s %-20s %12s\n", "ID", "Название", "Категория", "Цена")
	fmt.Println(tableLine)
	for _, p := range products {
		fmt.Printf("%-5d %-25s %-20s %12.2f\n", p.ID, p.Name, p.Category, p.Price)
	}
	fmt.Println(tableLine)
	fmt.Printf("Найдено товаров: %d\n", len(products))
}

// PrintSales печатает список продаж таблицей.
func PrintSales(sales []Sale) {
	if len(sales) == 0 {
		fmt.Println("Продажи не найдены.")
		return
	}
	fmt.Printf("%-10s %-25s %10s %-12s\n", "ProductID", "Название", "Кол-во", "Дата")
	fmt.Println(tableLine)
	for _, s := range sales {
		fmt.Printf("%-10d %-25s %10d %-12s\n", s.ProductID, s.Name, s.Quantity, s.Date)
	}
	fmt.Println(tableLine)
	fmt.Printf("Всего записей: %d\n", len(sales))
}

// PrintTotalRevenue печатает общую выручку и средний чек.
func PrintTotalRevenue(products []Product, sales []Sale) {
	total := TotalRevenue(products, sales)
	avg := 0.0
	if len(sales) > 0 {
		avg = total / float64(len(sales))
	}
	fmt.Println(tableLine)
	fmt.Printf("Общая выручка: %.2f\n", total)
	fmt.Printf("Средний чек:   %.2f  (по %d продажам)\n", avg, len(sales))
	fmt.Println(tableLine)
}

// PrintTopProducts печатает топ самых продаваемых товаров.
func PrintProductSales(items []ProductSales) {
	if len(items) == 0 {
		fmt.Println("Нет данных о продажах для топа.")
		return
	}
	fmt.Printf("%-4s %-5s %-25s %8s %14s\n", "#", "ID", "Название", "Кол-во", "Выручка")
	fmt.Println(tableLine)
	for i, ps := range items {
		fmt.Printf("%-4d %-5d %-25s %8d %14.2f\n",
			i+1, ps.Product.ID, ps.Product.Name, ps.Quantity, ps.Revenue)
	}
	fmt.Println(tableLine)
}

// PrintCategoryStats печатает статистику по категориям.
// Категории сортируются по выручке (по убыванию) — так информативнее.
func PrintCategoryStats(stats map[string]CategoryStat) {
	if len(stats) == 0 {
		fmt.Println("Статистика по категориям пуста.")
		return
	}

	names := make([]string, 0, len(stats))
	for name := range stats {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		// сначала по выручке, при равенстве — по имени
		if stats[names[i]].Revenue != stats[names[j]].Revenue {
			return stats[names[i]].Revenue > stats[names[j]].Revenue
		}
		return names[i] < names[j]
	})

	fmt.Printf("%-20s %10s %10s %14s\n", "Категория", "Товаров", "Продано", "Выручка")
	fmt.Println(tableLine)
	for _, name := range names {
		st := stats[name]
		fmt.Printf("%-20s %10d %10d %14.2f\n", name, st.Products, st.Quantity, st.Revenue)
	}
	fmt.Println(tableLine)
}

// PrintUnknownSales предупреждает о продажах, ссылающихся на несуществующие товары.
func PrintUnknownSales(sales []Sale) {
	if len(sales) == 0 {
		fmt.Println("Несогласованных продаж нет — все ID товаров найдены.")
		return
	}
	fmt.Printf("Внимание: %d продаж ссылаются на отсутствующие товары:\n", len(sales))
	PrintSales(sales)
}

// PrintProductsByCategory — обертка: фильтр + печать, чтобы не дублировать в main.
func PrintFilteredByCategory(products []Product, category string) {
	fmt.Printf("Товары в категории %q:\n", category)
	PrintProducts(FilterByCategory(products, category))
}

// PrintFilteredByPrice — обертка: фильтр по цене + печать.
func PrintFilteredByPrice(products []Product, min, max float64) {
	fmt.Printf("Товары с ценой от %.2f до %.2f:\n", min, max)
	PrintProducts(FilterByPrice(products, min, max))
}

// readLineTrim — вспомогательная функция для меню: читает строку и обрезает пробелы.
// Держим здесь, чтобы main.go не разрастался.
func readLineTrim(r *bufio.Reader) string {
	s, _ := r.ReadString('\n')
	return strings.TrimSpace(s)
}

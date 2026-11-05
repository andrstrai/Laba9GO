package main

import (
	"sort"
	"strings"
)

// ProductSales — сколько продано одного товара и сколько он принёс денег.
// Используется в топе самых продаваемых товаров.
type ProductSales struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
	Revenue  float64 `json:"revenue"`
}

// CategoryStat — статистика по одной категории товаров.
type CategoryStat struct {
	Products int     `json:"products"` // сколько товаров в категории
	Quantity int     `json:"quantity"` // сколько штук продано
	Revenue  float64 `json:"revenue"`  // выручка
}

// productIndex строит справочник «ID товара -> товар»,
// чтобы по продаже быстро найти товар, не перебирая весь срез.
func productIndex(products []Product) map[uint]Product {
	index := make(map[uint]Product, len(products))
	for _, p := range products {
		index[p.ID] = p
	}
	return index
}

// FilterByCategory возвращает товары нужной категории (регистр не важен).
func FilterByCategory(products []Product, category string) []Product {
	result := []Product{}
	for _, p := range products {
		if strings.EqualFold(p.Category, category) {
			result = append(result, p)
		}
	}
	return result
}

// FilterByPrice возвращает товары с ценой от min до max включительно.
func FilterByPrice(products []Product, min, max float64) []Product {
	result := []Product{}
	for _, p := range products {
		if p.Price >= min && p.Price <= max {
			result = append(result, p)
		}
	}
	return result
}

// TotalRevenue считает общую выручку: сумма цена * количество по всем продажам.
// Продажи товаров, которых нет в списке, не учитываются (см. UnknownSales).
func TotalRevenue(products []Product, sales []Sale) float64 {
	index := productIndex(products)
	total := 0.0
	for _, s := range sales {
		p, ok := index[s.ProductID]
		if !ok {
			continue
		}
		total += p.Price * float64(s.Quantity)
	}
	return total
}

// TopProducts возвращает n самых продаваемых товаров (по количеству штук).
// Если n <= 0 или n больше числа проданных товаров, возвращаются все.
func TopProducts(products []Product, sales []Sale, n int) []ProductSales {
	index := productIndex(products)

	// 1. Складываем продажи одного и того же товара.
	totals := make(map[uint]ProductSales)
	for _, s := range sales {
		p, ok := index[s.ProductID]
		if !ok {
			continue
		}
		ps := totals[s.ProductID]
		ps.Product = p
		ps.Quantity += s.Quantity
		ps.Revenue += p.Price * float64(s.Quantity)
		totals[s.ProductID] = ps
	}

	// 2. map не упорядочен, поэтому перекладываем в срез.
	result := make([]ProductSales, 0, len(totals))
	for _, ps := range totals {
		result = append(result, ps)
	}

	// 3. Сортируем: больше штук — выше; при равенстве — по ID.
	sort.Slice(result, func(i, j int) bool {
		if result[i].Quantity != result[j].Quantity {
			return result[i].Quantity > result[j].Quantity
		}
		return result[i].Product.ID < result[j].Product.ID
	})

	// 4. Обрезаем до n.
	if n > 0 && n < len(result) {
		result = result[:n]
	}
	return result
}

// CategoryStats считает по каждой категории: число товаров, проданные штуки и выручку.
func CategoryStats(products []Product, sales []Sale) map[string]CategoryStat {
	stats := make(map[string]CategoryStat)

	for _, p := range products {
		st := stats[p.Category]
		st.Products++
		stats[p.Category] = st
	}

	index := productIndex(products)
	for _, s := range sales {
		p, ok := index[s.ProductID]
		if !ok {
			continue
		}
		st := stats[p.Category]
		st.Quantity += s.Quantity
		st.Revenue += p.Price * float64(s.Quantity)
		stats[p.Category] = st
	}
	return stats
}

// UnknownSales возвращает продажи, у которых нет товара с таким ID.
// Нужна, чтобы предупредить пользователя о несогласованных файлах.
func UnknownSales(products []Product, sales []Sale) []Sale {
	index := productIndex(products)
	result := []Sale{}
	for _, s := range sales {
		if _, ok := index[s.ProductID]; !ok {
			result = append(result, s)
		}
	}
	return result
}

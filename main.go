package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type Product struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}
type Sale struct {
	ProductID uint   `json:"product_id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	Date      string `json:"date"`
}

func LoadProductsFromJSON(filename string, product *[]Product) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Ошибка чтения имени файла.")
	}
	var products []Product
	err = json.Unmarshal(data, &products)
	if err != nil {
		return fmt.Errorf("Ошибка разбора JSON.")
	}
	for i := 0; i < len(products); i++ {
		fmt.Println("----------")
		fmt.Println("ID: ", products[i].ID)
		fmt.Println("Name: ", products[i].Name)
		fmt.Println("Category: ", products[i].Category)
		fmt.Println("Price: ", products[i].Price)
		// cоздаем срез product и заполняем ее данными из JSON
		productss := Product{}
		productss.ID = products[i].ID
		productss.Name = products[i].Name
		productss.Category = products[i].Category
		productss.Price = products[i].Price
		*product = append(*product, productss)
	}
	return nil
}

func LoadSalesFromCSV(filename string, sales *[]Sale) error {
	data, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Ошибка чтения имени файла.")
	}
	defer data.Close()

	bufferedReader := bufio.NewReader(data)

	r := csv.NewReader(bufferedReader)
	r.Comma = ';'
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break // Файл закончился
		} else if err != nil {
			fmt.Println("Ошибка чтения файла.")
			break
		}
		if len(rec) < 4 {
			continue // Пропуск битой строки
		}
		fmt.Println("----------")
		fmt.Println("ID: ", rec[0])
		fmt.Println("Name: ", rec[1])
		fmt.Println("Quantity: ", rec[2])
		fmt.Println("Date: ", rec[3])
		// cоздаем срез Sale и заполняем ее данными из CSV
		sale := Sale{}
		id, err := strconv.ParseUint(rec[0], 10, 64)
		if err != nil {
			fmt.Println("Ошибка: ID товара не является числом в строке", rec)
			continue // Пропускаем строку
		}
		sale.ProductID = uint(id)
		sale.Name = rec[1]
		sale.Quantity, err = strconv.Atoi(rec[2])
		if err != nil {
			fmt.Println("Ошибка: количество не является числом в строке", rec)
			continue // Пропускаем строку
		}
		sale.Date = rec[3]
		*sales = append(*sales, sale)
	}
	return nil
}

// функция для очистки экрана (работает и на Windows, и на macOS/Linux)
func clear_screen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	// массивы для хранения продуктов и продаж
	myproducts := []Product{}
	mysales := []Sale{}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Для начала работы пожалуйста, нажмите Enter...")
		reader.ReadString('\n')
		clear_screen()

		fmt.Println("Меню команд:")
		fmt.Println("1 - Прочитать JSON файл с продуктами")
		fmt.Println("2 - Прочитать CSV файл с продажами")
		fmt.Println("3 - Сформировать отчёт report.json")
		fmt.Println("4 - Общая выручка")
		fmt.Println("5 - Топ продаваемых товаров")
		fmt.Println("6 - Фильтр товаров по категории")
		fmt.Println("7 - Фильтр товаров по цене")
		fmt.Println("8 - Статистика по категориям")
		fmt.Println("9 - Несогласованные продажи")
		fmt.Println("0 - Завершить работу")
		fmt.Println("Выберите пункт из меню управления: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		a, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Некорректный ввод! Пожалуйста, введите число.")
			continue
		}

		switch a {
		case 0:
			return
		case 1:
			fmt.Println("Пожалуйста, введите имя файла: ")
			input, _ := reader.ReadString('\n')
			filename := strings.TrimSpace(input)
			err := LoadProductsFromJSON(filename, &myproducts)
			if err != nil {
				fmt.Println(err)
			}
		case 2:
			fmt.Println("Пожалуйста, введите имя файла: ")
			input, _ := reader.ReadString('\n')
			filename := strings.TrimSpace(input)
			err := LoadSalesFromCSV(filename, &mysales)
			if err != nil {
				fmt.Println(err)
			}
		case 3:
			err := WriteReport(myproducts, mysales, 10, "report.json")
			if err != nil {
				fmt.Println("Ошибка:", err)
			}
		case 4:
			PrintTotalRevenue(myproducts, mysales)

		case 5:
			fmt.Print("Сколько товаров показать в топе (Enter — 10): ")
			nStr := readLineTrim(reader)
			n := 10
			if nStr != "" {
				var err error
				n, err = strconv.Atoi(nStr)
				if err != nil || n <= 0 {
					fmt.Println("Некорректное число, показываю топ-10.")
					n = 10
				}
			}
			PrintProductSales(TopProducts(myproducts, mysales, n))

		case 6:
			fmt.Print("Введите категорию: ")
			PrintFilteredByCategory(myproducts, readLineTrim(reader))

		case 7:
			fmt.Print("Минимальная цена: ")
			min, err := strconv.ParseFloat(readLineTrim(reader), 64)
			if err != nil {
				fmt.Println("Некорректная цена.")
				break
			}
			fmt.Print("Максимальная цена: ")
			max, err := strconv.ParseFloat(readLineTrim(reader), 64)
			if err != nil {
				fmt.Println("Некорректная цена.")
				break
			}
			if min > max {
				min, max = max, min // не падать, а просто поменять местами
			}
			PrintFilteredByPrice(myproducts, min, max)

		case 8:
			PrintCategoryStats(CategoryStats(myproducts, mysales))

		case 9:
			PrintUnknownSales(UnknownSales(myproducts, mysales))
		default:
			fmt.Println("Нет такого пункта меню!")
		}
	}
}

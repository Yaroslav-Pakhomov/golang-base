package interfaceWork

import (
	"fmt"
	"strings"
)

func GetAllStruct() {
	user := User{5, "Alex"}
	report := Report{"Title report", 7}
	invoice := Invoice{"№556-UTEH", 10}

	items := []Printer{user, &report, invoice} // ВАЖНО: &report

	PrintAll(items)
}

// Printer - Интерфейс
type Printer interface {
	Print() string
}

// =============================================================================================

// Несколько структур

type User struct {
	Id   int
	Name string
}

func (u User) Print() string {
	return fmt.Sprintf("Пользователь %v с ID %v печатает.", u.Name, u.Id)
}

type Report struct {
	Title string
	Page  uint
}

// Print - pointer receiver (r *Report): меняем состояние и показываем нюанс с реализацией интерфейса
func (r *Report) Print() string {
	// допустим, при печати автоматически добавляем пометку
	if !strings.HasPrefix(r.Title, "[PRINTED] ") {
		r.Title = "[PRINTED] " + r.Title
	}

	return fmt.Sprintf("Report: %s (%d pages)", r.Title, r.Page)
}

type Invoice struct {
	Number string
	Total  uint
}

func (i Invoice) Print() string {
	return fmt.Sprintf("Invoice Number: %v, Total: %v", i.Number, i.Total)
}

// =============================================================================================

// PrintAll - Функция, принимающая интерфейс
func PrintAll(items []Printer) {
	for idx, item := range items {
		fmt.Printf("%d) %s\n", idx+1, item.Print())
	}
}

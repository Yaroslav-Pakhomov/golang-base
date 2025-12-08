package main

import (
	"fmt"
	"golang-base/pkg/askname"
	"golang-base/pkg/greeting"
)

func main() {
	fmt.Println("Hello Go.")
	fmt.Println("My name is Иван.")
	fmt.Println("I’m starting to learn Go!")

	// Переменные и конкатенация строк
	greeting.Greeting()
	// Ввод данных от пользователя
	askname.Askname()
}

package varConst

import "fmt"

// 1. Переменные в Go
// var a int
// var b float64
// var name string

// var x, y int
// var width, height float64

// Краткая запись var a int = 10
// a := 10

// 2. Константы в Go
// const Pi = 3.14159
// const Greeting = "Hello, World!"

// 3. Типы данных
// Целые числа: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64.
// Числа с плавающей запятой: float32, float64.
// Строки: string.
// 	Логические значения: bool (имеет два значения: true и false).
// Массивы и срезы, структуры, карты (map) — типы данных, которые могут быть полезны для более сложных структур данных.

// 4. Инициализация переменных и констант
// var a int = 10
// const Pi = 3.14159
// "Краткая форма" объявления переменных
// a := 10

// 5. Важные особенности
// Переменные в Go всегда инициализируются значением по умолчанию, если явно не указано другое значение.
// Например, int по умолчанию равен 0, float64 — 0.0, string — пустая строка, а bool — false.

// 6. Операции с переменными
// Переменные могут быть использованы в различных операциях, таких как:
// Арифметические операции: +, -, *, /, %.
// Операции сравнения: ==, !=, <, >, <=, >=.
// Логические операции: && (и), || (или), ! (не).

func SetVars() {
	age := 25
	name := "Alice"
	isStudent := true
	height := 1.75
	numbers := [5]int{1, 4, 5, 3, 8}

	var strStudent string
	if isStudent {
		strStudent = "Да"
	} else {
		strStudent = "Нет"
	}

	fmt.Println("Возраст:", age)
	fmt.Println("Имя:", name)
	fmt.Println("Студент:", strStudent)
	fmt.Println("Рост:", height)
	fmt.Println("Номера:", numbers)
	fmt.Println("")
}

func SetConst() {
	const Pi = 3.14159
	const MAX_USERS = 100

	fmt.Println("Число Пи:", Pi)
	fmt.Println("Макс. кол-во пользователей:", MAX_USERS)
}

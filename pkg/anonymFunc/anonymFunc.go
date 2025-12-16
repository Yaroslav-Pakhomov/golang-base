package anonymFunc

import "fmt"

// Что такое анонимная функция
// Анонимная функция — это функция без имени, которую можно:
// 	- присвоить переменной,
// 	- передать как аргумент в другую функцию,
// 	- вызвать немедленно.

// В Go это выглядит так:
// func(параметры) (результаты) {
// 	// тело
// }

func GetAnonymFuncs() {

	// 1) Немедленный вызов (IIFE). Скобки () в конце — это вызов.
	func() {
		fmt.Println("Немедленный вызов")
	}()
	fmt.Println("")

	// 2) Присваивание переменной и вызов
	anonym := func() {
		fmt.Println("Анонимная функция в переменной")
	}
	anonym()
	fmt.Println("")

	// 3) Анонимная функция с параметрами
	func(params string) {
		fmt.Println(params)
	}("Анонимная функция с параметрами")
	fmt.Println("")

	// 4) Анонимная функция с возвращаемым значением
	resSumma := func(a, b int) int {
		return a + b
	}(2, 3)
	fmt.Println("Анонимная функция с возвращаемым значением сумма:", resSumma)
	fmt.Println("")

	// Замыкание-счётчик
	counter := 0
	next := func() int {
		counter++
		return counter
	}
	fmt.Println("Замыкание-счётчик вызов:", next())
	fmt.Println("Замыкание-счётчик вызов:", next())
	fmt.Println("Замыкание-счётчик вызов:", next())
	fmt.Println("")

	// Функция apply
	summa := apply(8, 4, func(x, y int) int { return x + y })
	fmt.Println("Сумма с помощью анонимной функции:", summa)
	multiple := apply(8, 4, func(x, y int) int { return x * y })
	fmt.Println("Произведение с помощью анонимной функции:", multiple)
	fmt.Println("")

	// Функция makePrefixer
	hello := makePrefixer("Hello, ")
	fmt.Println(hello("Go!"))
	fmt.Println("")

	// Функция makeCounter
	counterTest := makeCounter(1)
	fmt.Println("Счётчик шаг:", counterTest())
	counterTest = makeCounter(10)
	fmt.Println("Счётчик шаг:", counterTest())
}

// Напиши функцию apply(a, b int, op func(int, int) int) int, которая применяет op к a и b.
//   - Вызови apply(8, 4, func(x, y int) int { return x + y })
//   - И ещё раз для умножения и разности (используй анонимные функции).
func apply(a, b int, op func(int, int) int) int {
	result := op(a, b)
	return result
}

// Реализуй makePrefixer(prefix string) func(string) string, которая возвращает функцию, добавляющую prefix к входной строке.
// Пример: hello := makePrefixer("Hello, "), затем hello("Go") → "Hello, Go"
func makePrefixer(prefix string) func(string) string {
	return func(s string) string {
		return prefix + s
	}
}

// Реализуй makeCounter(step int) func() int:
//   - Внутри хранится счётчик, стартует с 0
//   - Каждый вызов увеличивает на step и возвращает текущее значение
//   - Проверь для step = 1 и step = 10
func makeCounter(step int) func() int {
	counter := 0

	return func() int {
		counter += step
		return counter
	}
}

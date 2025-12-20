package loop

import "fmt"

// Виды циклов:

// 1. Классический цикл for init; condition; post {},
// 		- init — инициализация (выполняется один раз перед циклом),
// 		- condition — условие продолжения (проверяется перед каждой итерацией),
// 		- post — шаг (выполняется после каждой итерации).

// 2. Условный цикл (аналог while),
// 		for condition {
// 			тело цикла
// 		}
// condition == true — цикл выполняется.

// 3. Бесконечный цикл,
// 		for {
// 			тело
// 		}

// 4. Цикл for range для перебора коллекций.
// 		for index, value := range collection {
// 			используем index и value
// 		}
// Если что-то не нужно (index, value) — можно использовать _.

// Используется для перебора:
// 		- срезов,
// 		- массивов,
// 		- строк,
// 		- карт (map),
// 		- каналов (в спец. случае).

func GetLoop() {
	// Классический
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}
	fmt.Println("")

	// Условный
	x := 1
	for x != 0 {
		fmt.Println("Введите число (0 для выхода): ")
		fmt.Scan(&x)
	}
	fmt.Println("Цикл завершён")
	fmt.Println("")

	// Бесконечный
	for {
		fmt.Println("1) Сказать привет")
		fmt.Println("2) Выйти")
		fmt.Print("Ваш выбор: ")

		var choice int
		fmt.Scan(&choice)

		if choice == 1 {
			fmt.Println("Привет!")
		} else if choice == 2 {
			fmt.Println("Выход...")
			fmt.Println("")
			break
		} else {
			fmt.Println("Неизвестный пункт меню")
		}
		fmt.Println("")
	}
	// fmt.Println("")

	// Диапазон range
	nums := []int{10, 20, 30, 40, 50}

	// range по срезу/массиву
	for index, value := range nums {
		fmt.Println("Индекс", index, " - значение", value)
	}
	fmt.Println("")

	for _, value := range nums {
		fmt.Println("Только значение", value)
	}
	fmt.Println("")

	// range по строке
	s := "Привет"
	// Итерирует по рунам (Unicode-символам), а не по байтам:
	for i, v := range s {
		fmt.Printf("позиция: %d, руна: %c\n", i, v)
	}
	fmt.Println("")

	// range по map
	ages := map[string]int{
		"Алиса": 25,
		"Вова":  30,
	}
	for name, age := range ages {
		fmt.Println(name, ":", age)
	}
	fmt.Println("")

	// Управляющие операторы: break и continue

	// Пример с continue: вывести только чётные числа от 0 до 9:
	for i := 0; i < 10; i++ {
		if i%2 != 0 {
			continue // пропускаем нечётные
		}
		fmt.Println(i)
	}
	fmt.Println("")

	// Пример с break: остановиться, когда найдём число > 100:
	numbers := []int{10, 20, 101, 30, 40, 50}
	for _, value := range numbers {
		if value > 100 {
			fmt.Println(value)
			break
		}
	}
	fmt.Println("")
}

func GetPracticeLoop() {
	var a, b int

	fmt.Println("Введите размеры таблицы умножения")

	fmt.Println("Первый параметр")
	fmt.Scan(&a)

	fmt.Println("Второй параметр")
	fmt.Scan(&b)

	for i := 1; i <= a; i++ {
		var multipl string
		for j := 1; j <= b; j++ {
			multipl += fmt.Sprintf("%d * %d = %d     ", i, j, i*j)
		}
		fmt.Println(multipl)
		fmt.Println("")
	}
}

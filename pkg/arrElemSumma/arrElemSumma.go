package arrElemSumma

import "fmt"

// Массивы и срезы в Go

// В Go существуют:
// Массивы — имеют фиксированную длину
// Срезы (slice) — динамические структуры (используются чаще)

// Пример массива:
// var arr [5]int = [5]int{1, 2, 3, 4, 5}

// Пример среза:
// numbers := []int{1, 2, 3, 4, 5}

// ===================================================

// Цикл for в Go

// В Go используется один универсальный цикл for.
// Перебор элементов массива или среза:

// for i := 0; i < len(numbers); i++ {
// // работа с numbers[i]
// }

// Или более удобный вариант с range:
// for _, value := range numbers {
// // value — текущий элемент
// }

func GetArrayElementSumma() {
	numbers := []int{4, 6, 2, 10, 5}

	summa := summaArray(numbers)
	fmt.Println("Сумма элементов массива:", summa)

	SetArrayFor()
}

func summaArray(arr []int) (summa int) {

	for _, number := range arr {
		summa += number
	}

	return
}

// Задача 1: посчитать сумму массива произвольной длины и значений:
// 	- пользователь вводил количество элементов массива
// 	- пользователь вводил сами элементы
// 	- программа выводила сумму элементов

func SetArrayFor() {

	// Ввод размера массива
	var arrLen int
	fmt.Println("Введите кол-во элементов массива:")
	_, err := fmt.Scan(&arrLen)
	if err != nil || arrLen <= 0 {
		fmt.Println("Ошибка: введите корректное положительное число")
		SetArrayFor()
		return
	}

	// Создание среза нужного размера
	numbers := make([]int, arrLen)

	// Ввод элементов массива
	for i := 0; i < arrLen; i++ {
		fmt.Println("Введите эл-т с индексом:", i)
		_, err := fmt.Scan(&numbers[i])
		if err != nil {
			fmt.Println("Ошибка: введено некорректное значение")
			return
		}
	}

	// Подсчёт суммы
	summa := summaArray(numbers)

	fmt.Println("Сумма элементов массива:", summa)
}

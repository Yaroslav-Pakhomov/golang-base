package switchCase

import "fmt"

// GetDay - базовый switch
func GetDay(day string) {

	switch day {
	case "Mon":
		fmt.Println("Сегодня пн.")
	default:
		fmt.Println("Неизвестный день.")
	}

}

// GetGraduate - Несколько значений в одном case
func GetGraduate(grade int) {

	switch grade {
	case 5, 4:
		fmt.Println("Отлично / Хорошо")
	case 3:
		fmt.Println("Удовлетворительно")
	case 2, 1, 0:
		fmt.Println("Неудовлетворительно")
	default:
		fmt.Println("Некорректная оценка")
	}
}

// GetTemperature - switch без выражения (условный)
func GetTemperature(temp int) {

	switch {
	case temp < 0:
		fmt.Println("Мороз")
	case temp >= 0 && temp < 10:
		fmt.Println("Холодно")
	case temp >= 10 && temp < 20:
		fmt.Println("Тепло")
	default:
		fmt.Println("Жарко")
	}
}

// CheckX - Если нужно «провалиться» в следующий case — есть fallthrough
func CheckX(x int) {

	switch x {
	case 1:
		fmt.Println("Первая")
		fallthrough
	case 2:
		fmt.Println("Вторая (сразу после первой)")
	default:
		fmt.Println("Остальные")
	}
}

// Задание 1. Классификация символов
// 	Дан символ ch типа rune. Вывести, к какой группе он относится:
// 	 - гласная: a, e, i, o, u (и в верхнем регистре тоже)
// 	 - согласная латиница (достаточно обработать несколько примеров: b, c, d, f, g и т.п.)
// 	 - цифра: 0..9
// 	 - другое
// 	Подсказка: используй switch и несколько значений в одном case.
// 	Подсказка: rune в Go — это символ Unicode.

// GetRune - Классификация символов
func GetRune(ch rune) {

	switch ch {
	// гласные
	case 'a', 'e', 'i', 'o', 'u',
		'A', 'E', 'I', 'O', 'U':
		fmt.Println("Гласная")

	// согласные
	case 'b', 'c', 'd', 'f', 'g',
		'B', 'C', 'D', 'F', 'G':
		fmt.Println("Согласная")

	// цифры
	case '0', '1', '2', '3', '4',
		'5', '6', '7', '8', '9':
		fmt.Println("Цифра")

	// всё остальное
	default:
		fmt.Println("другое")
	}
}

package ifElse

import "fmt"

// GetEven - узнать чётность/нечётность числа
func GetEven() {
	var a int

	fmt.Println("Введите число \"a\":")
	fmt.Scan(&a)
	checkEven(a)

	var b int

	fmt.Println("Введите число \"b\":")
	fmt.Scan(&b)
	checkEven(b)

	var c int

	fmt.Println("Введите число \"c\":")
	fmt.Scan(&c)
	checkEven(c)
}

func checkEven(number int) {
	if number == 0 {
		fmt.Println("Ноль — это особое число")
	} else if number%2 == 0 {
		fmt.Println("Число чётное")
	} else {
		fmt.Println("Число нечётное")
	}

	if number > 0 {
		fmt.Println("Положительное")
	}

	if number < 0 {
		fmt.Println("Отрицательное")
	}
}

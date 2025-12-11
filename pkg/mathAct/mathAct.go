package mathAct

import "fmt"

// Summa - рассчёт суммы трёх чисел
func Summa() {
	var a, b, c int

	fmt.Println("Введите первое число:")
	fmt.Scan(&a)

	fmt.Println("Введите второе число:")
	fmt.Scan(&b)

	fmt.Println("Введите третье число:")
	fmt.Scan(&c)

	sum := a + b + c
	fmt.Println("Сумма чисел:", sum)
}

// Calculator - выполняет действия: +, -, *, /
func Calculator() {
	var a, b int
	var operation string

	fmt.Println("Введите первое число:")
	fmt.Scan(&a)

	fmt.Println("Введите второе число:")
	fmt.Scan(&b)

	fmt.Println("Какую операцию выполнить: +, -, *, /")
	fmt.Scan(&operation)

	switch operation {
	case "+":
		fmt.Println("Сумма чисел:", a+b)
	case "-":
		fmt.Println("Разность между первым и вторым:", a-b)
	case "*":
		fmt.Println("Произведение чисел:", a*b)
	case "/":
		if b != 0 {
			fmt.Println("Деление первого на второе:", a/b)
		} else {
			fmt.Println("Деление на 0")
		}
	}
}

// ArithmeticMean - Среднее арифметическое 4-ёх чисел
func ArithmeticMean() {
	var a, b, c, d int

	fmt.Println("Введите четыре числа:")
	fmt.Scan(&a)
	fmt.Scan(&b)
	fmt.Scan(&c)
	fmt.Scan(&d)

	fmt.Println("Среднее арифметическое 4-ёх чисел равно:", (a+b+c+d)/4)
}

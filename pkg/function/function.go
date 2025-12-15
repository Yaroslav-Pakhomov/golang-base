package function

import (
	"fmt"
)

// func имяФункции(параметры) (тип1, тип2) {
// 	return значение1, значение2
// }

func RequestFunctions() {
	a := 4
	b := 1
	summa, diff := SummaAndDiff(a, b)
	fmt.Println("Сумма ", a, " и ", b, "равна", summa, ", разность", a, " и ", b, "равна", diff)

	// Если не нужно одно из возвращаемых значений, используется _ (пустой идентификатор):
	s, _ := SummaAndDiff(a, b)
	fmt.Println("Сумма ", a, " и ", b, "равна", s)

	a = 10
	b = 2
	div, err := Divide(a, b)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Деление ", a, " на ", b, "равна", div)
	}

	a = 10
	b = 0
	div, err = Divide(a, b)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Деление ", a, " на ", b, "равна", div)
	}

	a = 5
	b = 10
	sq, perimeter := RectangleParams(a, b)
	fmt.Println("Площадь прямоугольника со сторонами ", a, " и ", b, "равна", sq, ", периметр - ", perimeter)

}

// SummaAndDiff - считает сумму и разность чисел
func SummaAndDiff(a int, b int) (sum int, diff int) {
	sum = a + b
	diff = a - b

	return sum, diff
}

// Divide - паттерн: результат + ошибка
func Divide(a int, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("Деление на 0.")
		// errors.New("divide by zero")
	}

	return a / b, nil
}

// RectangleParams - считает площадь и периметр прямоугольника
func RectangleParams(x int, y int) (sq int, perimeter int) {
	sq = x * y
	perimeter = (x + y) * 2

	return sq, perimeter
}

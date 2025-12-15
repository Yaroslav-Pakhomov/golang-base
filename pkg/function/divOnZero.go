package function

import (
	"errors"
	"fmt"
)

// Сигнатура функции с ошибкой:
// func someFunc(...) (resultType, error)

// Для создания ошибок используется стандартный пакет:
// import "errors"

// Пример создания ошибки:

// errors.New("описание ошибки")

// В Go принято не допускать панику (panic) во время выполнения, а явно возвращать ошибку.

func CheckFunc() {
	result, err := DivideWithError(10, 2)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Результат:", result)

	result, err = DivideWithError(10, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Результат:", result)
}

func DivideWithError(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("Деление на 0 запрещено.")
	}

	return a / b, nil
}

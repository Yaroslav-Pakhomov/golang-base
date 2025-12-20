package funcNamedVal

import (
	"errors"
	"fmt"
	"strconv"
)

// Когда это полезно
// Именованные результаты полезны, когда:
// 	 - Возвращается несколько значений, и имена сразу объясняют смысл (written, n, ok, err, status, result).
// 	 - Хочется централизовать возврат и обработку ошибок (но не превращать функцию в «лабиринт»).
// 	 - Нужно, чтобы defer мог модифицировать возвращаемый результат (частый случай: логирование/метрики/оборачивание ошибки).

// defer в Go — это ключевое слово, которое откладывает выполнение функции до момента, когда текущая функция завершает работу.
// Проще:
// defer = «сделай это перед выходом из функции».

// Когда лучше НЕ использовать
// Избегайте, если:
// 	- Функция длинная, много ветвлений: naked return начинает скрывать, что именно возвращается.
// 	- Имена не помогают (типа (a int, b int) без смысла).
// 	- Возвращаемые значения легко выразить явно — иногда лучше написать return x, nil, чем держать в голове состояние переменных.

func GetFuncsWithNamedVals() {

	q, r, err := divide(4, 3)
	if err == nil {
		fmt.Println("Деление целая часть:", q, ", Остаток от деление:", r)
	} else {
		fmt.Println("Ошибка:", err)
	}

	q, r, err = divide(4, 0)
	if err == nil {
		fmt.Println(q, r)
	} else {
		fmt.Println("Ошибка:", err)
	}
	fmt.Println("")

	port, err := parsePort("8080")
	if err == nil {
		fmt.Println(port)
	} else {
		fmt.Println(err)
	}

	port, err = parsePort("abc")
	if err == nil {
		fmt.Println(port)
	} else {
		fmt.Println(port, err)
	}

	port, err = parsePort("70000")
	if err == nil {
		fmt.Println(port)
	} else {
		fmt.Println(port, err)
	}
}

// Особенности:
//   - q, r, err становятся переменными функции.
//   - Можно писать return без значений — вернутся текущие значения этих переменных (это называется naked return).
func divide(a, b int) (q, r int, err error) {
	if b == 0 {
		return 0, 0, errors.New("Деление на 0 запрещено.")
	}

	q = a / b
	r = a % b

	return
}

// Задание 1. Рефакторинг «безымянных» возвратов
func parsePort(s string) (port int, err error) {
	port, err = strconv.Atoi(s)

	if err != nil {
		return 80, errors.New("error")
	}

	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("invalid port: %d", port)
	}

	return port, nil
}

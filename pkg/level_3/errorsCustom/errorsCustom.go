package errorsCustom

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Пользовательская errors

// region Теоретическая часть

// Что такое error в Go

// В Go ошибка — это значение, реализующее интерфейс:
// 	type error interface {
// 		Error() string
// 	}

// То есть любая структура (или тип), у которой есть метод Error() string, считается ошибкой.

// -----------------

// Зачем делать свой тип ошибки

// Свой тип ошибки полезен, когда нужно:
// 	- отличать “свои” ошибки от чужих (например, ошибка бизнес-логики vs ошибка сети);
// 	- хранить дополнительные поля: код, параметры, имя файла, операцию;
// 	- удобно определять тип ошибки при обработке (errors.As).

// -----------------

// Оборачивание ошибок (wrapping)

// Go-стиль: не терять исходную ошибку, а добавлять контекст:
// return fmt.Errorf("read config: %w", err)

// Так можно потом проверить первопричину через errors.Is/As.

// -----------------

// errors.Is и errors.As

// 	- errors.Is(err, target) — проверяет, есть ли в цепочке обёрток конкретная ошибка (например, os.ErrNotExist).
// 	- errors.As(err, &targetType) — пытается “привести” ошибку к конкретному типу (например, *MyError).

// -----------------

// Обработка нескольких ошибок в одной функции

// Частый сценарий: одна функция вызывает несколько операций (I/O, парсинг, валидация), и нужно разрулить разные ошибки по-разному:
// 	- если файл не найден — один ответ;
// 	- если невалидные данные — другой;
// 	- если любая иная ошибка — третий.

// endregion Теоретическая часть

// region Практическая часть

// Задание 1. Создаём собственный тип ошибки MyError

// Требования:
// 	- создать type MyError struct{ ... }
// 	- добавить поля: Op string, Code int, Msg string
// 	- реализовать Error() string

// -----------------

// Задание 2. Одна функция — несколько возможных ошибок

// Реализуйте функцию LoadAgeFromFile(path string) (int, error):
// 	- читает файл,
// 	- парсит число (возраст),
// 	- валидирует диапазон (например, 0…150),
// 	- возвращает разные ошибки:
// 		- если файл не существует — вернуть обёрнутую системную ошибку с контекстом;
// 		- если число не парсится — вернуть MyError с кодом (например, 400);
// 		- если возраст вне диапазона — вернуть MyError с другим кодом (например, 422).

// -----------------

// Задание 3. Расширяем MyError

// 	- Добавьте в MyError поле Field string (например, имя поля данных).
// 	- Измените Error() так, чтобы при наличии Field оно выводилось в сообщении.
// Ожидаемо: ошибки парсинга/валидации указывают, какое поле неверное (например, Field="age").

func GetTestCustomErrors() {
	age, err := LoadAgeFromFile("age.txt")
	if err != nil {
		HandleError(err)
	}

	fmt.Println("Возраст:", age)

	fmt.Println("")
	PrintUserParseResult("Ivan;25") // ok
	fmt.Println("")
	PrintUserParseResult("Ivan;999") // 422
	fmt.Println("")
	PrintUserParseResult("Ivan;abc") // 400 age
	fmt.Println("")
	PrintUserParseResult(";25") // 422 name
	fmt.Println("")
	PrintUserParseResult("Ivan25") // 400 line
}

type MyError struct {
	Op    string
	Code  int
	Msg   string
	Field string
}

func (e *MyError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("op = %s code = %d field = %s msg = %s", e.Op, e.Code, e.Field, e.Msg)
	}
	return fmt.Sprintf("op = %s, code = %d, message = %s", e.Op, e.Code, e.Msg)
}

// LoadAgeFromFile - читает возраст из файла, парсит и валидирует.
func LoadAgeFromFile(path string) (int, error) {
	const op = "LoadAgeFromFile"

	data, errRead := os.ReadFile(path)
	if errRead != nil {
		// добавляем контекст, не теряя первопричину
		return 0, fmt.Errorf("%s: чтение файла: %w", op, errRead)
	}

	strData := strings.TrimSpace(string(data))
	age, errConv := strconv.Atoi(strData)
	if errConv != nil {
		// это уже ошибка данных (не системная), можно вернуть свой тип
		return 0, &MyError{
			Op:    op,
			Code:  400,
			Msg:   fmt.Sprintf("невозможно распарсить возраст из %q", strData),
			Field: "age",
		}
	}

	if age < 0 || age > 150 {
		// это уже ошибка данных (не системная), можно вернуть свой тип
		return 0, &MyError{
			Op:    op,
			Code:  422,
			Msg:   fmt.Sprintf("возраст д.б. в диапазоне от 0 до 150, сейчас %d", age),
			Field: "age",
		}
	}

	return age, nil
}

// HandleError - показывает обработку нескольких ошибок в одной функции.
func HandleError(err error) {
	if err == nil {
		fmt.Println("Ok")
		return
	}
	// 1) Проверка на "файл не найден" через errors.Is
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Ошибка: файл не найден")
		fmt.Println("Технические детали:", err)
		return
	}

	// 2) Извлечение MyError через errors.As
	var myError *MyError
	if errors.As(err, &myError) {
		switch myError.Code {
		case 400:
			fmt.Println("Ошибка ввода данных (не число).", myError)
		case 422:
			fmt.Println("Ошибка ввода данных (вне диапазона).", myError)
		default:
			fmt.Println("Моя ошибка.", myError)
		}
		return // <-- важно
	}

	// 3) Всё остальное
	fmt.Println("Неизвестная ошибка:", err)
}

// -----------------

// Задание 4. Функция с несколькими источниками ошибок

// Реализуйте функцию:
// 	- func ParseUser(line string) (name string, age int, err error)

// Формат строки: "Name;Age", например "Alice;21".
// Требования:
// 	- Если нет ";" или частей меньше/больше двух — вернуть MyError (код 400, Field="line").
// 	- Если возраст не число — MyError (код 400, Field="age").
// 	- Если возраст вне диапазона 0..150 — MyError (код 422, Field="age").
// 	- Если имя пустое — MyError (код 422, Field="name").

func ParseUser(line string) (name string, age int, err error) {
	const op = "ParseUser"

	lineParts := strings.SplitN(line, ";", 2)
	if len(lineParts) != 2 {
		// Свой тип ошибки
		return "", 0, &MyError{
			Op:    op,
			Code:  400,
			Msg:   fmt.Sprintf("невозможно распарсить строку для получения имени и возраста из %q, ожидается \"Name;Age\"", line),
			Field: "line",
		}
	}

	// Имя
	name = strings.TrimSpace(lineParts[0])
	if name == "" {
		return "", 0, &MyError{
			Op:    op,
			Code:  422,
			Msg:   fmt.Sprintf("некорректно имя из %q", line),
			Field: "name",
		}
	}

	// Возраст
	ageStr := strings.TrimSpace(lineParts[1])
	age, err = strconv.Atoi(ageStr)
	if err != nil {
		// Свой тип ошибки
		return "", 0, &MyError{
			Op:    op,
			Code:  400,
			Msg:   fmt.Sprintf("некорректный возраст из %q", ageStr),
			Field: "age",
		}
	}

	if age < 0 || age > 150 {
		// Свой тип ошибки
		return "", 0, &MyError{
			Op:    op,
			Code:  422,
			Msg:   fmt.Sprintf("возраст д.б. в диапазоне от 0 до 150, сейчас %d", age),
			Field: "age",
		}
	}

	return name, age, nil
}

// -----------------

// Задание 5. Единая обработка ошибок

// Напишите функцию:
// 	- func PrintUserParseResult(line string)

// Она вызывает ParseUser и в одной функции обрабатывает ошибки:
// 	- Code=400 → “Неверный формат входных данных”
// 	- Code=422 → “Данные корректны по формату, но не проходят валидацию”
// 	- иначе → “Неизвестная ошибка”

// Используйте errors.As.

func PrintUserParseResult(line string) {
	nameLine, ageLine, errLine := ParseUser(line)

	if errLine == nil {
		fmt.Printf("Имя: %s, возраст: %d.\n", nameLine, ageLine)
		return
	}

	var myError *MyError
	if errors.As(errLine, &myError) {
		switch myError.Code {
		case 400:
			fmt.Println("Неверный формат входных данных (строка не корректна, проверьте наличие \";\").", myError)
		case 422:
			fmt.Println("Данные корректны по формату, но не проходят валидацию (не корректно имя и/или возраст).", myError)
		default:
			fmt.Println("Неизвестная ошибка", myError)
		}
		return
	}

	fmt.Println("Неизвестная ошибка:", errLine)
}

// endregion Практическая часть

// region Контрольные вопросы

// 1. Что такое интерфейс error в Go и какие требования к типу, чтобы он считался ошибкой?
// error — это интерфейс с одним методом:
// 		type error interface { Error() string }
// Любой тип (структура, алиас), у которого есть метод Error() string, считается ошибкой и может возвращаться как error.

// 2. Зачем “оборачивать” ошибку через %w, а не просто форматировать %v?
// %w сохраняет исходную ошибку внутри “обёртки”, формируя цепочку причин. Тогда можно:
// 	- находить первопричину через errors.Is(err, os.ErrNotExist),
// 	- извлекать тип через errors.As(err, &myErr).
// Если сделать %v, это будет просто строка — цепочки причин не будет.

// 3. Разница между errors.Is и errors.As?
// 	- errors.Is(err, target) — проверяет “есть ли в цепочке” конкретная ошибка/маркер (например os.ErrNotExist).
// 	- errors.As(err, &targetType) — пытается привести ошибку к нужному типу (например *MyError) и получить доступ к полям.

// 4. Почему полезно добавлять поле Op (операция) в пользовательскую ошибку?
// Op показывает, где произошла ошибка (имя функции/операции). Это ускоряет отладку и логирование: по ошибке видно контекст без поиска по коду.

// 5. Как удобнее организовать обработку нескольких ошибок и почему?
// Часто удобно так:
// 	- сначала errors.Is для “системных маркеров” (файл не найден и т.п.),
// 	- затем errors.As + switch по коду/полю для своих ошибок,
// 	- в конце “всё остальное”.
// switch по Code читается лучше длинной if-цепочки. Таблица соответствий (map) удобна, когда много кодов и одинаковый обработчик.

// 6. Когда лучше возвращать системную ошибку, а когда — свою (MyError)?
// 	- Системные (I/O, сеть, права, не существует файл) — лучше не терять и возвращать как причину, добавив контекст через wrapping (%w).
// 	- Бизнес/валидация/формат данных — лучше возвращать свою (MyError), чтобы иметь коды, поле, понятные сообщения и стабильную обработку.

// endregion Контрольные вопросы

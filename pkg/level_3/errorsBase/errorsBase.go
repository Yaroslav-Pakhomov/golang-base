package errorsBase

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// Основное определение errors

// region Теоретическая часть

// Интерфейс error

// В Go ошибка — это значение, реализующее интерфейс:
// 	type error interface {
// 		Error() string
// 	}

// Типичный стиль: функция возвращает (результат, error):
// 	- если всё ок → err == nil
// 	- если проблема → err != nil, а результат часто нулевой (0, "", nil, пустая структура — зависит от типа)

// -----------------

// Создание ошибок через errors.New()

// errors.New("текст") создаёт простую ошибку с сообщением:
// 	err := errors.New("file not found")

// Когда удобно:
// 	- для простых, статичных ошибок (без параметров)
// 	- как “сигнальную” ошибку, которую потом можно сравнивать через errors.Is()

// -----------------

// Оборачивание ошибок через fmt.Errorf("%w", err)

// Оборачивание нужно, чтобы:
// 	- добавить контекст (где и что сломалось),
// 	- но не потерять исходную причину.

// Пример:
// return fmt.Errorf("read config: %w", err)

// Важно:
// 	-  %w используется только для error и только один раз в строке форматирования.
// 	-  Если использовать %v вместо %w, то ошибка будет “склеена” текстом, и errors.Is() уже не найдёт оригинал.

// -----------------

// Проверка обёрнутых ошибок
// 	- errors.Is(err, target) — ищет target внутри цепочки обёрток
// 	- errors.Unwrap(err) — возвращает “внутреннюю” ошибку (если она была обёрнута)

// -----------------

// Когда нужен errors.As
// Функция	Когда использовать
// errors.Is	проверить конкретную “сигнальную” ошибку (ErrEmptyInput, os.ErrNotExist)
// errors.As	достать ошибку определённого типа и получить доступ к её полям

// errors.As не сравнивает, а извлекает.

// endregion Теоретическая часть

// region Практическая часть

var (
	ErrEmptyInput    = errors.New("пустая строка")
	ErrAgeOutOfRange = errors.New("возраст вне допустимого диапазона")
	ErrDivideByZero  = errors.New("деление на 0")
	ErrEmptyPath     = errors.New("пустой путь к директории")
	ErrInvalidConfig = errors.New("invalid config")
)

func GetTestParseAge() {

	// ParseAge
	tests := []string{"", "abc", "-5", "200", "42"}

	for _, test := range tests {
		age, err := ParseAge(test)

		if err != nil {
			fmt.Printf("input=%q -> error: %v\n", test, err)

			// Примеры проверок:
			if errors.Is(err, ErrEmptyInput) {
				fmt.Println("  причина: пустой ввод")
			}
			if errors.Is(err, ErrAgeOutOfRange) {
				fmt.Println("  причина: возраст вне диапазона")
			}
			continue
		}

		fmt.Printf("input=%q -> age=%v\n", test, age)
	}

	// SafeDivide
	testsDiv := [][2]int{
		{10, 5},
		{7, 3},
		{5, 0},
		{15, 3},
	}

	for _, item := range testsDiv {
		result, err := SafeDivide(item[0], item[1])
		if err != nil {
			fmt.Printf("%d / %d: %v\n", item[0], item[1], err)

			if errors.Is(err, ErrDivideByZero) {
				fmt.Println(" причина: деление на ноль")
			}
			continue
		}

		fmt.Printf("%d / %d: %v\n", item[0], item[1], result)
	}

	// LoadFile
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	_, err := LoadFile("")
	fmt.Println("LoadFile(\"\"):", err)

	_, err = LoadFile("no_such_file.txt")
	fmt.Println("LoadFile(\"no_such_file.txt\"):", err)

	if err != nil && errors.Is(err, os.ErrNotExist) {
		fmt.Println("  причина: файла не существует (errors.Is сработал)")
	}

	// initApp
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	errInit := initApp("no_such_file.txt")
	if errInit != nil {
		fmt.Println("error:", errInit)
		if errors.Is(errInit, os.ErrNotExist) {
			fmt.Println("  errors.Is(err, os.ErrNotExist) == true")
		}
	}

}

// Задание 1. Реализуй функцию, которая возвращает error

// Требование: внутри использовать:
// 	- errors.New() для базовой ошибки
// 	- fmt.Errorf("%w") для оборачивания

// Условие
// Напиши функцию ParseAge(input string) (int, error):
// 	- Если input пустая строка → вернуть ошибку ErrEmptyInput (создать через errors.New).
// 	- Попробовать преобразовать строку в число.
// 	- Если преобразование не удалось → вернуть ошибку, обёрнутую через fmt.Errorf("invalid age: %w", err).
// 	- Если возраст меньше 0 или больше 150 → вернуть ErrAgeOutOfRange (через errors.New).
// 	- Иначе вернуть возраст и nil.

// ParseAge - Парсинг и проверка возраста
func ParseAge(input string) (int, error) {
	if input == "" {
		return 0, ErrEmptyInput
	}

	age, err := strconv.Atoi(input)
	if err != nil {
		// strconv.Atoi возвращает *strconv.NumError.
		var numError *strconv.NumError
		if errors.As(err, &numError) {
			fmt.Println("")
			fmt.Println("ошибка преобразования числа")
			fmt.Println(" функция:", numError.Func)
			fmt.Println(" значение:", numError.Num)
			fmt.Println(" причина:", numError.Err)
			fmt.Println("")
		}
		return 0, fmt.Errorf("invalid age: %w", err)
	}

	if age < 0 || age > 150 {
		return 0, ErrAgeOutOfRange
	}

	return age, nil
}

// -----------------

// Задание 2(база). “Калькулятор деления”

// Требования:
// 	- Если b == 0 → вернуть ErrDivideByZero (через errors.New).
// 	- Иначе вернуть результат и nil.

// SafeDivide - Задание 2(база): “Калькулятор деления”
func SafeDivide(a, b int) (int, error) {

	if b == 0 {
		return 0, ErrDivideByZero
	}

	return a / b, nil
}

// -----------------

// Задание 3(с оборачиванием). “Чтение файла с контекстом”

// Требования:
// 	- Если path == "" → вернуть ErrEmptyPath (через errors.New).
// 	- Иначе попытаться прочитать файл os.ReadFile(path).
// 	- Ошибку чтения обернуть: fmt.Errorf("load file %q: %w", path, err).

func LoadFile(path string) ([]byte, error) {

	if path == "" {
		return nil, ErrEmptyPath
	}

	readFile, err := os.ReadFile(path)
	if err != nil {

		var pathError *os.PathError
		if errors.As(err, &pathError) {
			fmt.Println("")
			fmt.Println("ошибка преобразования файла")
			fmt.Println(" операция", pathError.Op)
			fmt.Println(" путь", pathError.Path)
			fmt.Println(" причина", pathError.Err)
			fmt.Println("")
		}

		return nil, fmt.Errorf("load file %q: %w", path, err)
	}

	return readFile, nil
}

// -----------------

// Задание 4(усложнение). “Слои приложения”

// - `readConfig(path string) ([]byte, error)` — читает файл, оборачивает ошибку.
// - `parseConfig(data []byte) (Config, error)` — если пусто, возвращает `ErrInvalidConfig`.
// - `initApp(path string) error` — вызывает обе, и каждую ошибку оборачивает своим контекстом.

// Цель: получить цепочку ошибок вида:
// - `init app: read config: load file "x": <оригинальная ошибка>`

// И показать:
// - что `errors.Is(err, os.ErrNotExist)` срабатывает,
// - а сообщение ошибки остаётся понятным человеку.

// Config - структура для хранения содержания файла
type Config struct {
	Raw string
}

// initApp - вызывает обе, и каждую ошибку оборачивает своим контекстом
func initApp(path string) error {
	data, err := ReadConfig(path)
	if err != nil {
		return fmt.Errorf("read file %q: %w", path, err)
	}

	_, err = ParseConfig(data)
	if err != nil {
		return fmt.Errorf("parse file %q: %w", path, err)
	}

	return nil
}

// ParseConfig - парсинг содержания файла
func ParseConfig(data []byte) (Config, error) {
	if len(data) == 0 {
		return Config{}, ErrInvalidConfig
	}

	return Config{Raw: string(data)}, nil
}

// ReadConfig - Чтение файла с контекстом
func ReadConfig(path string) ([]byte, error) {
	b, err := loadFile(path)

	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return b, nil
}

// loadFile - Загрузка файла с контекстом
func loadFile(path string) ([]byte, error) {
	if path == "" {
		return nil, ErrEmptyPath
	}

	readFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load file %q: %w", path, err)
	}

	return readFile, nil
}

// endregion Практическая часть

// region Контрольные вопросы — ответы

// Почему возвращают ошибку отдельным значением, а не исключения?
// Потому что Go предпочитает явный контроль потока: ошибки — часть обычного результата функции. Это делает поведение предсказуемым, упрощает понимание и не требует “магических” try/catch.

// Чем отличается errors.New() от fmt.Errorf()?
// errors.New() — создаёт простую статичную ошибку (только текст).
// fmt.Errorf() — создаёт ошибку с форматированием, и может оборачивать другую ошибку через %w.

// Разница между %v и %w в fmt.Errorf?
// %v просто вставляет текст ошибки (цепочка теряется).
// %w создаёт “обёртку” и сохраняет исходную ошибку внутри — тогда её можно найти через errors.Is() / достать через errors.Unwrap().

// Что делает errors.Is() и когда полезен?
// Проверяет, есть ли target в цепочке обёрнутых ошибок. Полезен, когда вы добавляете контекст, но хотите различать причины (например, os.ErrNotExist, context.Canceled, свои сигнальные ошибки).

// Можно ли использовать %w два раза в одной строке? Почему?
// Нельзя. fmt.Errorf допускает ровно один %w, потому что ошибка может содержать одну “внутреннюю” причину (одну цепочку). Для нескольких причин используют другие подходы (например, errors.Join в новых Go).

// Что возвращает errors.Unwrap(err) и когда вернёт nil?
// Возвращает “внутреннюю” ошибку, если err её хранит (обёртка).
// Вернёт nil, если ошибка не обёрнута (или err == nil).

// Почему важно добавлять контекст при пробросе ошибок наверх?
// Без контекста наверху остаётся “file not found”, и непонятно: какой файл? какой шаг? какая функция? Контекст превращает ошибку в понятный трейс действий, при этом %w сохраняет причину для errors.Is.

// endregion Контрольные вопросы — ответы

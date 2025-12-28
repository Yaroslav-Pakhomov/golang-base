package emptyInterface

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Пустой Интерфейс

// region Теоретическая часть

// Пустой интерфейс interface{}

// В Go интерфейс — это набор методов.
// Пустой интерфейс interface{} не требует ни одного метода, поэтому любой тип удовлетворяет ему автоматически.

// Пример: функция принимает что угодно:
// 	func PrintAny(v interface{}) {
// 		fmt.Println(v)
// 	}

// -----------------

// Когда используют interface{}
// 	- когда нужно принять значение неизвестного заранее типа;
// 	- когда работаете с API, где тип данных динамический (например, JSON как map[string]interface{}).

// -----------------

// Риски и минусы
// 	- теряется статическая типизация на границе функции;
// 	- можно легко ошибиться с ожидаемым типом и получить панику при неверном приведении (если не проверять ok).

// -----------------

// Как вернуть конкретный тип: type assertion
// 	s, ok := v.(string)
// 	if !ok {
// 		// v не строка
// 	}

// -----------------

// Type switch (удобно для нескольких вариантов)
// 	switch x := v.(type) {
// 	case int:
// 		fmt.Println("int:", x)
// 	case string:
// 		fmt.Println("string:", x)
// 	default:
// 		fmt.Println("unknown type")
// 	}

// !!! В современном Go часто предпочитают generics (обобщения), но interface{} всё ещё встречается и важно понимать, как с ним работать.

// ===================================================================================

// Методы структуры: value receiver vs pointer receiver

// Пусть есть структура:
// 	type Counter struct {
// 		Value int
// 	}

// -----------------

// Value receiver (получатель-значение)
// 	func (c Counter) Inc() {
// 		c.Value++ // изменяется копия
// 	}
// 	c — копия структуры. Снаружи Value не изменится.

// -----------------

// Pointer receiver (получатель-указатель)
// 	func (c *Counter) Inc() {
// 		c.Value++ // изменяем оригинал через указатель
// 	}

// Теперь метод меняет исходный объект.

// Когда нужен pointer receiver
// 	- метод изменяет поля структуры;
// 	- структура тяжёлая (копировать дорого);
// 	- нужно единое поведение и совместимость с интерфейсами (иногда метод-сет у T и *T различается).

// endregion Теоретическая часть

// region Практическая часть

// GetWorkCheck - проверка работы Заданий
func GetWorkCheck() {

	// Задание 1 - Универсальная функция на interface{}
	Describe(5)
	Describe(true)
	Describe("Строка №5")
	Describe(64.5)
	Describe(5 + 64.5i)
	fmt.Println("")
	fmt.Println("")

	// Задание 2 - Структура и метод с указателем
	bAcc := &BankAccount{"Yar", 500}

	fmt.Println(bAcc.Deposit(-2))
	fmt.Println(bAcc.Deposit(20))

	fmt.Println(bAcc.Withdraw(-1000))
	fmt.Println(bAcc.Withdraw(1000))
	fmt.Println(bAcc.Withdraw(300))

	fmt.Println(bAcc.String())

	// Задание 3 - совместить Задания 1 и 2
	_ = Apply("deposit", bAcc, "50")
	_ = Apply("withdraw", bAcc, 10)
	fmt.Println(bAcc.String())
}

// region Задание 1: универсальная функция на interface{}

// Реализуйте функцию Describe(v interface{}) string, которая возвращает строку-описание значения:
// 	- если int → "целое число: <значение>"
// 	- если float64 → "вещественное число: <значение>"
// 	- если string → "строка: <значение>"
// 	- если bool → "логическое: <значение>"
// 	- иначе → "неизвестный тип"

// Подсказка: используйте type switch.

func Describe(v interface{}) {

	switch t := v.(type) {
	case int:
		fmt.Printf("Целое число: %v. \n", t)
	case float64:
		fmt.Printf("Вещественное число: %v. \n", t)
	case string:
		fmt.Printf("Строка: %v. \n", t)
	case bool:
		fmt.Printf("Логическое: %v. \n", t)
	default:
		fmt.Printf("Неизвестный тип: %v. \n", t)
	}
}

// endregion Задание 1

// -----------------

// region Задание 2: структура и метод с указателем

// Создайте структуру BankAccount:
// 	- Owner string
// 	- Balance int

// Реализуйте методы:

// func (a *BankAccount) Deposit(amount int) error
// Увеличивает баланс. Ошибка, если amount <= 0.

// func (a *BankAccount) Withdraw(amount int) error
// Уменьшает баланс. Ошибка, если amount <= 0 или недостаточно средств.

// func (a BankAccount) String() string
// Возвращает строку вида: "Owner=<...>, Balance=<...>"

// Обратите внимание:
// 	- Deposit и Withdraw должны быть с pointer receiver, иначе изменения не сохранятся.
// 	- String() можно сделать value receiver — он не изменяет объект.

// BankAccount - Банковский личный кабинет
type BankAccount struct {
	Owner   string
	Balance float64
}

// Deposit - Увеличивает баланс
func (ba *BankAccount) Deposit(amount float64) string {
	if amount < 0 {
		return fmt.Sprintf("Недостаточно средств: %v. \n", amount)
	}

	ba.Balance += amount

	return fmt.Sprintf("Баланс пополнен: %v. \n", ba.Balance)
}

// Withdraw - Уменьшает баланс
func (ba *BankAccount) Withdraw(amount float64) string {
	if amount < 0 {
		return fmt.Sprintf("Недостаточно средств: %v. \n", amount)
	}

	if ba.Balance < amount {
		return fmt.Sprintf("Недостаточно средств на балансе: %v. \n", ba.Balance)
	}

	ba.Balance -= amount

	return fmt.Sprintf("Баланс уменьшился на %v: %v. \n", amount, ba.Balance)
}

// String - Возвращает строку с информацией о Пользователе
func (ba *BankAccount) String() string {
	return fmt.Sprintf("Пользователь: %v, Баланс: %v. \n", ba.Owner, ba.Balance)
}

// endregion Задание 2

// -----------------

// region Задание 3: совместить оба навыка

// Напишите функцию:
// func Apply(op string, account *BankAccount, value interface{}) error

// Требования:
// 	- op может быть "deposit" или "withdraw";
// 	- value может быть:
// 	- int (использовать как сумму),
// 	- float64 (округлить вниз до int),
// 	- string (попробовать распарсить в число),
// 	- иначе → ошибка "unsupported type".

// После извлечения суммы вызвать соответствующий метод у account.

// Apply - выполняет операцию deposit/withdraw, а сумму получает из interface{}.
// Поддерживаем int, float64 (floor), string (парсинг).
func Apply(op string, account *BankAccount, value interface{}) error {

	if account == nil {
		return errors.New("аккаунт не существует")
	}

	amount, err := extractValue(value)
	if err != nil {
		return err
	}

	switch op {
	case "deposit":
		account.Deposit(amount)
	case "withdraw":
		account.Withdraw(amount)
	default:
		_ = errors.New("неизвестная операция")
	}

	return nil
}

// extractAmount — выделим логику извлечения суммы в отдельную функцию, чтобы Apply был коротким и понятным.
func extractValue(value interface{}) (float64, error) {

	switch t := value.(type) {

	case int:
		return float64(t), nil

	case float64:
		return t, nil

	case string:
		s := strings.TrimSpace(t)
		if s == "" {
			return 0, errors.New("empty string amount")
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("невозможно распарсить значение из строки %q: %w", t, err)
		}
		return float64(n), nil

	case bool:
		return 0, nil

	default:
		return 0, fmt.Errorf("некорректное значение")
	}
}

// endregion Задание 3

// endregion Практическая часть

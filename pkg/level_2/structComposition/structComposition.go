package structComposition

import "fmt"

// Композиция структур

// region Теоретическая часть

// 1. Почему не наследование?

// В Go нет классического наследования классов. Вместо этого используется:
// 	- композиция (struct внутри struct),
// 	- встраивание (embedding) — частный случай композиции, когда тип встраивается без имени поля.

// Идея: “не быть кем-то” (is-a), а “содержать” (has-a).

// ===================================================================================

// 2. Обычное поле vs встраивание

// Обычное поле:

// 	type Engine struct {
// 		Power int
// 	}

// 	type Car struct {
// 		Eng Engine // поле с именем Eng
// 	}

// Доступ: car.Eng.Power

// -----------------

// Встраивание:

// 	type Car struct {
// 		Engine // встроенный тип (без имени поля)
// 	}

// Доступ: car.Power (поле “поднимается” наверх), а также car.Engine.Power тоже работает.

// ===================================================================================

// 3. Продвижение (promotion) полей и методов

// Если тип T встроен в S, то:
// 	- поля T доступны как будто они объявлены в S,
// 	- методы T также доступны через S.

// Это удобно для “сборки” поведения:
// 	- S получает методы T “из коробки”,
// 	- но при этом остаётся независимым типом.

// ===================================================================================

// 4. Конфликты имён

// Если во встроенных типах встречаются одинаковые имена:
// 	- прямой доступ s.Name станет неоднозначным → ошибка компиляции,
// 	- нужно явно указать путь: s.Type1.Name или s.Type2.Name.

// ===================================================================================

// 5. Переопределение методов (shadowing)

// Если внешний тип определит метод с тем же именем, что и встроенный:
// 	- вызов s.Method() использует метод внешнего типа,
// 	- но встроенный метод доступен по полному пути: s.EmbeddedType.Method().

// ===================================================================================

// 6. Встраивание интерфейсов

// Интерфейсы можно встраивать друг в друга:

// type Reader interface { Read(p []byte) (int, error) }
// type Closer interface { Close() error }

// type ReadCloser interface {
// 	Reader
// 	Closer
// }

// ===================================================================================

// endregion Теоретическая часть

// region Практическое часть

// Практическое задание (основное)
// Задание: “Сервис доставки — композиция без наследования”
// Реализуй небольшую модель, где “собирается” объект заказа из компонентов.

// Требования:

// 1. Создай структуры:

// Address:
// 	- City string
// 	- Street string
// 	- House string

// Customer:
// 	- Name string
// 	- Phone string
// 	- встроенный Address

// Item:
// 	- Name string
// 	- Price int (в копейках/центах, чтобы без float)
// 	- Qty int

// Order:
// 	- ID string
// 	- встроенный Customer
// 	- Items []Item

// -----------------

// 2. Добавь методы:

// Для Item:
// 	- Cost() int → Price * Qty

// Для Order:
// 	- Total() int → сумма Cost() по всем товарам
// 	- AddItem(item Item) → добавляет в список
// 	- Summary() string → строка с краткой информацией (ID, клиент, адрес, total)

// 3. В main():

// Создай заказ,
// добавь 2–3 товара,
// выведи:
// 	- имя и телефон клиента через “поднятые” поля,
// 	- город/улицу/дом через “поднятые” поля (из Address),
// 	- итоговую сумму,
// 	- Summary().

// Ожидаемый акцент: ты должен использовать именно встраивание Address в Customer и Customer в Order, чтобы обращаться к полям как order.City, order.Name и т.п.

// endregion Практическое часть

// region Решение

func GetMainWork() {

	// region Создание заказа
	customer := Customer{
		"John Doe",
		"777",
	}

	item := Item{
		Name:  "Soap",
		Price: 222.88,
		Qty:   5,
	}

	items := []Item{item}

	order := Order{
		"Uteh-Or-01",
		customer,
		items,
	}
	// endregion Создание заказа

	// region Добавление 2–3 товаров
	item1 := Item{
		Name:  "apple",
		Price: 555,
		Qty:   3,
	}

	item2 := Item{
		Name:  "banana",
		Price: 33,
		Qty:   6,
	}

	order.AddItem(item1)
	order.AddItem(item2)

	// endregion Добавление 2–3 товаров

	// region Вывод инфы по заказу

	// Имя и телефон клиента через “поднятые” поля
	fmt.Printf("Пользователь %v, телефон %v. \n", customer.Name, customer.Phone)

	// Город/улицу/дом через “поднятые” поля (из Address)
	address := Address{
		"Мск",
		"Ленина",
		"55 корп. А",
		customer,
	}
	fmt.Printf("Город %v, улица %v дом %v. \n", address.City, address.Street, address.House)

	// Итоговая сумма
	fmt.Printf("Итоговая сумма заказа %v р. \n", order.Total())

	// Summary()
	fmt.Println(order.Summary())

	// endregion Вывод инфы по заказу
}

// Address - адрес Заказчика
type Address struct {
	City, Street, House string
	Customer
}

// Customer - Заказчик
type Customer struct {
	Name, Phone string
}

// Item - единица заказа
type Item struct {
	Name  string
	Price float64
	Qty   int
}

// Cost - стоимость единицы
func (item *Item) Cost() float64 {
	return item.Price * float64(item.Qty)
}

// Order - заказ
type Order struct {
	ID string
	Customer
	Items []Item
}

// Total - сумма заказа
func (order *Order) Total() float64 {
	var total float64

	for _, item := range order.Items {
		total += item.Cost()
	}

	return total
}

// AddItem - добавление единицы заказа
func (order *Order) AddItem(item Item) {
	// Item{
	// 	Name:  "Soap",
	// 	Price: 222.88,
	// 	Qty:   5,
	// }

	order.Items = append(order.Items, item)
}

// Summary - общая инфа по заказу
func (order *Order) Summary() string {
	return fmt.Sprintf("Пользователь %v, телефон %v с ID заказа %v.", order.Name, order.Phone, order.ID)
}

// endregion Решение

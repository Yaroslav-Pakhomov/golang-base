package structJson

import (
	"encoding/json"
	"fmt"
)

// Структура с JSON-тегами и Сериализация и десериализация структуры в JSON

// region Теоретическая часть

// Что такое JSON-теги в Go

// В Go можно «подсказать» пакету encoding/json, как поле структуры должно называться в JSON и как оно должно себя вести при кодировании/декодировании.

// Тег пишется после типа поля в обратных кавычках:

// 	type User struct {
// 		ID   int    `json:"id"`
// 		Name string `json:"name"`
// 	}

// 	- json:"id" — имя ключа в JSON будет id, а не ID.
// 	- Без тегов encoding/json использует имя поля как есть (ID, Name) и применяет правила экспорта (поле должно начинаться с заглавной буквы, иначе его не будет в JSON).

// =============================================================================================

// Основные правила

// Экспортируемость: поле должно быть экспортируемым (Name, а не name), иначе encoding/json его проигнорирует.

// Опции тега:
// 	omitempty — поле не попадёт в JSON, если значение «пустое» (0, "", nil, false, пустой срез/мапа).
// 	- — поле полностью игнорируется при маршалинге и анмаршалинге.

// 	type Profile struct {
// 		Email    string `json:"email,omitempty"`
// 		Password string `json:"-"` // никогда не попадёт в JSON
// 	}

// =============================================================================================

// Сериализация и десериализация

// 	- json.Marshal(v) → []byte с JSON
// 	- json.Unmarshal(data, &v) → заполнение структуры из JSON

// b, _ := json.Marshal(user)
// _ = json.Unmarshal(b, &user2)

// =============================================================================================

// Вложенные структуры, срезы и указатели

// 	- Вложенная структура кодируется «внутрь» JSON-объекта.
// 	- Срез ([]T) превращается в JSON-массив.
// 	- Указатель (*T) удобно использовать вместе с omitempty, чтобы отличать «нет значения» от «нулевого значения».

// endregion Теоретическая часть

// region Практическая часть

// Задание: описать модель заказа и поработать с JSON

// Шаг 1. Создай структуру Order с JSON-тегами

// Требования к JSON-формату:
// 	- order_id (число)
// 	- customer (объект)
// 	- items (массив объектов)
// 	- comment (строка, опционально, не включать если пустая)
// 	- internal_note — поле существует в структуре, но не должно попадать в JSON

// Шаг 2. Добавь вложенные структуры

// Customer: name, phone
// Item: sku, qty, price

// Шаг 3. Сериализация

// 	- Создай пример заказа в коде.
// 	- Преобразуй его в JSON (красиво, с отступами).
// 	- Выведи JSON в консоль.

// Шаг 4. Десериализация

// 	- Возьми строку JSON (пример ниже).
// 	- Распарсь её в структуру Order.
// 	- Выведи результат так, чтобы было видно, что поля прочитались.

// Пример JSON для Unmarshal:

// 	{
// 		"order_id": 1001,
// 		"customer": {
// 			"name": "Ivan Petrov",
// 			"phone": "+7-900-000-00-00"
// 		},
// 		"items": [
// 			{ "sku": "A-100", "qty": 2, "price": 199.90 },
// 			{ "sku": "B-200", "qty": 1, "price": 499.00 }
// 		],
// 		"comment": ""
// 	}

// GetWorkStructs - работа со структурой и JSON
func GetWorkStructs() {
	customer := Customer{"YAR", "777"}
	items := []Item{
		{
			"Apple",
			"55",
			35.5,
		},
		{
			"Banana",
			"35",
			55.8,
		},
		{
			"Peach",
			"15",
			255.3,
		},
	}

	// Marshal
	order := Order{1, customer, items, "Быстро", "Пожалуйста"}

	pretty, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("MARSHAL результат:")
	fmt.Println(string(pretty))

	// Unmarshal
	inputOrder := []byte(`{
	  "order_id": 1001,
	  "customer": {
		"name": "Ivan Petrov",
		"phone": "+7-900-000-00-00"
	  },
	  "items": [
		{ "sku": "A-100", "qty": "2", "price": 199.90 },
		{ "sku": "B-200", "qty": "1", "price": 499.00 }
	  ],
	  "comment": ""
	}`)

	var order2 Order
	if errUnmarsh := json.Unmarshal(inputOrder, &order2); errUnmarsh != nil {
		panic(errUnmarsh)
	}
	fmt.Println("\nUNMARSHAL результат:")
	fmt.Printf("%+v \n", order2)
}

type Customer struct {
	Name  string
	Phone string
}

type Item struct {
	Sku   string
	Qty   string
	Price float64
}

type Order struct {
	OrderId      int `json:"order_id"`
	Customer     `json:"customer"`
	Items        []Item `json:"items"`
	Comment      string `json:"comment,omitempty"`
	InternalNote string `json:"-"`
}

// endregion Практическая часть

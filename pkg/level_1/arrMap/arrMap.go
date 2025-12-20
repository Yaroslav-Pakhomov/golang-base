package arrMap

import (
	"fmt"
	"sort"
	"strings"
)

// Теоретическая часть

// Что такое map
// map — это ассоциативный массив (словарь): структура данных “ключ → значение”.

// Сигнатура:
// map[KeyType]ValueType

// Примеры:
// 	- map[string]int — имя → количество
// 	- map[int]string — id → название

// ===================================================================================

// Создание карты

// 1.Через make (самый частый вариант):
// m := make(map[string]int)

// 2. Литералом (удобно для стартовых данных):
// m := map[string]int{"apple": 3, "banana": 5}

// Важно: nil-map читать можно, но записывать нельзя.
// var m map[string]int // nil

// ===================================================================================

// Операции с map

// Запись/обновление:
// m["apple"] = 10

// Чтение:
// v := m["apple"] // если ключа нет, вернётся значение по умолчанию (0 для int)

// Проверка наличия ключа (“comma ok”):
// v, ok := m["apple"]
// if ok { /* ключ есть */ } else { /* ключа нет */ }

// Удаление:
// delete(m, "apple")

// Длина:
// n := len(m)

// ===================================================================================

// Проход по map

// for k, v := range m {
// fmt.Println(k, v)
// }

// Важный момент: порядок обхода map в Go не гарантирован и может меняться между запусками.

// ===================================================================================

// Как сделать вывод в отсортированном порядке

// Чтобы получить стабильный порядок, обычно:
// 	- собирают ключи в []T,
// 	- сортируют,
// 	- выводят значения по ключам.

func WorkWithMap() {
	GetBaseWorkWithMap()
	fmt.Println("")

	PracticeBaseWorkWithMap()
	fmt.Println("")

	GetAnaliseWord("Go go GO is fun fun")
}

func GetBaseWorkWithMap() {
	// Создание map
	store := make(map[string]int)

	// Добавление/обновление
	store["apple"] = 100
	store["peach"] = 250
	store["grape"] = 390
	store["banana"] = 200
	store["cherry"] = 300
	store["banana"] = 500 // обновили значение

	// Чтение + проверка наличия ключа
	if v, ok := store["apple"]; ok {
		fmt.Println("Яблоки в наличии", v)
	} else {
		fmt.Println("Яблок нет в наличии", v)
	}

	if v, ok := store["pear"]; ok {
		fmt.Println("Груша в наличии", v)
	} else {
		fmt.Println("Груш нет в наличии", v)
	}

	// Обход map (порядок НЕ гарантирован)
	fmt.Println("\nНеупорядоченная итерация:")
	for k, v := range store {
		fmt.Printf("%s: %d\n", k, v)
	}

	// Удаление
	delete(store, "cherry")

	// Стабильный вывод: сортируем ключи
	keys := make([]string, 0, len(store))
	for k := range store {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("\nУпорядоченная итерация:")
	for _, k := range keys {
		fmt.Printf("%s: %d\n", k, store[k])
	}

	// Длинна
	fmt.Println("\nОбщее количество товаров", len(store))
}

// Базовая задача:

// Создай map[string]int с оценками студентов (например: "Ivan": 5, "Anna": 4).
// Добавь ещё 3 записи, одну запись обнови.
// Выведи:
// 	- всех студентов и их оценки (через range),
// 	- количество студентов (len),
// 	- есть ли студент "Petr" (через v, ok := ...).

// PracticeBaseWorkWithMap - Базовая задача
func PracticeBaseWorkWithMap() {
	student := map[string]int{"Ivan": 5, "Anna": 4, "Jack": 3, "Ann": 2}

	student["Jacks"] = 4
	student["Liza"] = 3
	student["Elena"] = 3

	fmt.Println(student, "\n")

	for name, estimate := range student {
		fmt.Printf("Студент %s: оценка \"%d\".\n", name, estimate)
	}

	fmt.Println("\nОбщее кол-во студентов:", len(student))

	if petr, ok := student["Petr"]; ok {
		fmt.Println("Да, студент", petr, "есть.")
	} else {
		fmt.Println("Нет, студента Петра нет.")
	}
}

// Задача: Частотный словарь слов.

// Дан текст (строка). Нужно:
// 	- привести к нижнему регистру,
// 	- разбить на слова (можно считать, что слова разделены пробелами — без сложной пунктуации),
// 	- посчитать, сколько раз встречается каждое слово (map[string]int),
// 	- вывести слова и частоты в алфавитном порядке.

// GetAnaliseWord - Частотный словарь слов
func GetAnaliseWord(text string) {

	// Нижний регистр
	text = strings.ToLower(text)

	fmt.Println(text)
	fmt.Println("")

	// Разбиваем на слова (по пробелам/любым пробельным символам)
	words := strings.Fields(text)

	// Считаем число вхождений слова
	freq := make(map[string]int)
	for _, word := range words {
		freq[word]++
	}

	// Вывод в алфавитном порядке: собираем и сортируем ключи
	keys := make([]string, 0, len(freq))
	for k := range freq {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	// Печать результата
	for _, word := range keys {
		fmt.Printf("Слово %s: число вхождений -  %d\n", word, freq[word])
	}
}

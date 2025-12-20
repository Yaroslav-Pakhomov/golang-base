package structAndMeth

import "fmt"

// Теоретическая часть

// Структуры в Golang
// Структура (struct) — это пользовательский тип данных, который позволяет объединять несколько переменных (полей) под одним именем.

// Пример структуры:
// type Person struct {
// 	Name string
// 	Age  int
// }

// Здесь:
// 	- Person — имя структуры
// 	- Name, Age — поля структуры
// 	- string, int — типы данных полей

// ===================================================================================

// Создание экземпляра структуры

// Экземпляр структуры можно создать следующим образом:
// 		p := Person{
// 			Name: "Иван",
// 			Age:  25,
// 		}

// ===================================================================================

// Методы структуры

// Метод — это функция, связанная с определённым типом (структурой).

// Общий вид метода:
// 		func (receiver TypeName) MethodName() {
//			// тело метода
// 		}

// receiver — получатель (экземпляр структуры, к которому относится метод).

// В языке Go:
// 	- value receiver (Person) → метод получает копию
// 	- pointer receiver (*Person) → метод работает с оригиналом

// ===================================================================================

// GetStructs - общий метод по работе со структурами
func GetStructs() {

	// Создание экземпляра структуры
	person := Person{"Иван", 28, "Екатеринбург"}
	fmt.Println(person)

	// Вывод данных структуры
	fmt.Println("Имя:", person.Name)
	fmt.Println("Возраст:", person.Age)
	fmt.Println("")

	// Вызов метода
	person.Greeting()
	fmt.Println("")

	person.HaveBirthday()
	fmt.Println("Возраст", person.Age)
	fmt.Println("")

	if person.IsAdult() {
		fmt.Println("Возраст больше 18")
	} else {
		fmt.Println("Возраст меньше 18")
	}
}

// Практика:
// 1. Создать структуру Person
// 2. Добавить поля:
// 		- Name (строка)
// 		- Age (целое число)
// 3. Реализовать метод Greet(), который выводит приветствие
// 4. Вывести данные структуры в функции main

// Person - Структура пользователя
type Person struct {
	Name string
	Age  int
	City string
}

// Greeting - Представление пользователя со всеми данными
func (person *Person) Greeting() {
	fmt.Println("Привет! Меня зовут", person.Name, ", мой возраст", person.Age, "лет. Я живу в", person.City, ".")
}

// HaveBirthday - увеличение возраста. !!! Важно использовать pointer receiver (*Person)
func (person *Person) HaveBirthday() {
	person.Age++
}

// IsAdult - проверка на совершеннолетия
func (person *Person) IsAdult() bool {
	return person.Age >= 18
}

package structAndMethDeep

import (
	"fmt"
	"math"
)

func GetStructDeep() {

	// Прямоугольник
	var rectangle Shape
	rectangle = Rectangle{10.5, 10.5}

	fmt.Println("Площадь прямоугольника:", rectangle.Area())
	fmt.Println("Периметр прямоугольника:", rectangle.Perimeter())

	fmt.Println("")

	// Круг
	var circle Shape
	circle = Circle{10}

	fmt.Println("Площадь круга:", circle.Area())
	fmt.Println("Периметр круга:", circle.Perimeter())

	fmt.Println("")

	// Треугольник
	var triangle Shape
	triangle = Triangle{3, 4, 5}

	fmt.Println("Периметр треугольника:", triangle.Perimeter())
	fmt.Println("Площадь треугольника:", triangle.Area())

	fmt.Println("")

	// Срез со всеми фигурами
	shapes := []Shape{
		Circle{5},
		Rectangle{5.5, 5.5},
		Triangle{4.4, 5, 2},
	}

	for i, shape := range shapes {
		fmt.Printf("%d.", i+1)
		DescribeShape(shape)
		PrintShapesInfo(shape)
	}

	// Type Assertion
	PrintCircleRadius(rectangle)
	PrintCircleRadius(circle)
	fmt.Println("")

	// Type Switch
	DescribeShape(rectangle)
	DescribeShape(circle)
	DescribeShape(triangle)
}

// ===================================================================================

// PrintShapesInfo - функцию, которая выводит тип фигуры, её площадь и периметр
func PrintShapesInfo(s Shape) {
	fmt.Printf("Фигура (%T)\n", s)
	fmt.Printf("  Площадь: %.2f\n", s.Area())
	fmt.Printf("  Периметр: %.2f\n", s.Perimeter())
	fmt.Printf("-------------------------------------------\n\n")
}

// ===================================================================================

// Функция на применение Type Assertion

// PrintCircleRadius - проверка, что shape - это структура Circle
func PrintCircleRadius(shape Shape) {

	// Type Assertion
	checkCircle, success := shape.(Circle)

	// fmt.Println(checkCircle, success)
	if success {
		fmt.Printf("Радиус круга равен %v\n", checkCircle.radius)
	} else {
		fmt.Println("Фигура не является кругом")
	}
}

// Функция на применение Type Switch

// DescribeShape - проверка типа фигуры и вывод соответствующих данных
func DescribeShape(shape Shape) {
	switch shapeType := shape.(type) {
	case Rectangle:
		fmt.Printf("Фигура Прямоугольник. Ширина равна %v , высота %v. \n", shapeType.Width, shapeType.Height)
	case Circle:
		fmt.Printf("Фигура Круг. Радиус круга равен %v. \n", shapeType.radius)
	case Triangle:
		fmt.Printf("Фигура Треугольник. Сторона \"А\" равна %v, сторона \"В\" равна %v, сторона \"С\" равна %v. \n", shapeType.A, shapeType.B, shapeType.C)
	default:
		fmt.Printf("Неизвестная фигура")
	}
}

// ===================================================================================

// Shape - интерфейс Shape, который описывает геометрическую фигуру.
// Любая фигура, у которой есть методы Area() и Perimeter(), считается Shape.
type Shape interface {
	Area() float64
	Perimeter() float64
}

// ===================================================================================

// Здесь:
// 	- Rectangle — имя структуры
// 	- Width, Height — поля структуры
// 	- float64 — тип данных полей

// Rectangle - структура прямоугольника
type Rectangle struct {
	Width  float64
	Height float64
}

// Метод — это функция, связанная с конкретным типом (структурой).

// Area - Площадь прямоугольника
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter - Периметр прямоугольника
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// ===================================================================================

// Circle - структура круга
type Circle struct {
	radius float64
}

// Area - Площадь круга
func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

// Perimeter - Периметр круга
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

// ===================================================================================

// Triangle - структура треугольника
type Triangle struct {
	A float64
	B float64
	C float64
}

// Area - Площадь треугольника, формула Герона
func (t Triangle) Area() float64 {
	if !t.IsValid() {
		return 0 // или math.NaN()
	}

	p := t.Perimeter() / 2
	return math.Sqrt(p * (p - t.A) * (p - t.B) * (p - t.C))
}

// Perimeter - Периметр треугольника
func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

// IsValid - валидность для формулы Герона, сумма любых двух сторон должна быть строго больше третьей.
func (t Triangle) IsValid() bool {
	return t.A+t.B > t.C && t.A+t.C > t.B && t.B+t.C > t.A
}

// ===================================================================================

// Интерфейс
// type Speaker interface {
// 	Speak() string
// }

// Структура
// type Human struct {}
// func (h Human) Speak() string {
// 	return "Hello"
// }

// Type assertion — это проверка и извлечение конкретного типа из интерфейса.

// Синтаксис:
// value, ok := interfaceValue.(ConcreteType)

// 		- value — значение конкретного типа
// 		- ok — true, если приведение успешно
// 		- interfaceValue - интерфейс, который реализует структура
// 		- ConcreteType - структура, проверка на принадлежность

// Пример Type assertion:

// Переменная Структуры "Human", реализующая Интерфейс "Speaker"
// var s Speaker = Human{}

// Type assertion:
// 	h, ok := s.(Human)
// 	if ok {
// 		fmt.Println("Это Human:", h)
// 	}

// Опасно:
// h := s.(Human) // panic, если тип не Human

// ===================================================================================

// Type switch — удобный способ проверить несколько типов интерфейса.

// Синтаксис
// 	switch v := interfaceValue.(type) {
// 	case Type1:
// 		...
// 	case Type2:
// 		...
// 	default:
// 		...
// 	}

// Пример Type Switch со структурами:

// func describe(s Speaker) {
// 	switch v := s.(type) {
// 	case Human:
// 		fmt.Println("Human говорит:", v.Speak())
// 	case Dog:
// 		fmt.Println("Dog говорит:", v.Speak())
// 	default:
// 		fmt.Println("Неизвестный тип")
// 	}
// }

// Когда использовать:
// Ситуация							Решение
// Проверка одного типа				type assertion
// Много возможных типов			type switch
// Бизнес-логика зависит от типа	type switch
// Полиморфизм без проверок			интерфейсные методы

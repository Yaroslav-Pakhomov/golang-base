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
		PrintShapesInfo(shape)
	}
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

package structAndMethDeep

import (
	"fmt"
	"math"
)

func GetStructDeep() {

	// Прямоугольник
	rectangle := Rectangle{10.5, 10.5}

	fmt.Println("Площадь прямоугольника:", rectangle.Area())
	fmt.Println("Периметр прямоугольника:", rectangle.Perimeter())

	fmt.Println("")

	// Круг
	circle := Circle{10}

	fmt.Println("Площадь круга:", circle.Area())
	fmt.Println("Периметр круга:", circle.Perimeter())
}

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

package askname

import "fmt"

func Askname() {
	var name string

	fmt.Println("Введите ваше имя:")
	fmt.Scanln(&name)
	fmt.Println("Привет, ", name, "! Добро пожаловать в Go!")
}

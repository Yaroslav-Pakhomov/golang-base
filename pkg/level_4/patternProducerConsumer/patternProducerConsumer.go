package patternProducerConsumer

import (
	"fmt"
	"sync"
)

// Producer / Consumer (производитель / потребитель)

// region Теоретическая часть

// Producer / Consumer (производитель / потребитель)
// одни горутины производят данные (producer)
// другие обрабатывают (consumer)

// Пример простой с пояснениями:

// Классический Producer / Consumer в Go на горутинах и каналах — максимально просто и наглядно.

// Идея:
// 	- Producer (производитель) — генерирует данные и кладёт их в канал
// 	- Consumer (потребитель) — читает данные из канала и обрабатывает
// 	- Channel — очередь между ними
// 	- Goroutines — работают параллельно

// endregion Теоретическая часть

// region Практическая часть

func GetPatternProducerConsumer() {

	// 1-ый способ — sync.WaitGroup
	{
		// блок видимости - это отдельный кусок кода, со своими локальными переменными
		chInt := make(chan int)

		var wgInt sync.WaitGroup
		wgInt.Add(1)

		go producer(chInt)
		go func() {
			defer wgInt.Done()
			consumer(chInt)
		}()

		wgInt.Wait()
		fmt.Println("Работа с каналом закончена через WaitGroup.")
	}

	fmt.Println("")

	// 2-ой способ — канал done
	{
		chInt := make(chan int)
		done := make(chan struct{})

		go producer(chInt)
		go func() {
			consumer(chInt)
			close(done)
		}()

		<-done
		fmt.Println("Работа с каналом закончена через канал done.")
	}
}

// producer — генерирует данные и кладёт их в канал
func producer(chInt chan<- int) {
	for i := 0; i < 10; i++ {
		chInt <- i
	}

	close(chInt)
}

// consumer - читает данные из канала и обрабатывает
func consumer(chInt <-chan int) {
	for v := range chInt {
		fmt.Println(v)
	}
}

// region Разбор примера

// Что тут происходит (по шагам)
// 	- make(chan int)
// 		* создаём канал для передачи int
// 	- go producer(ch)
// 		* producer начинает параллельно генерировать числа
// 	- ch <- i
// 		* producer кладёт данные в канал
// 		* если consumer не успевает — producer ждёт
// 	- for value := range ch
// 		* consumer читает данные
// 		* цикл автоматически завершается, когда канал закрыт
// 	- close(ch)
// 		* producer сообщает: «данные закончились»

// endregion Разбор примера

// endregion Практическая часть

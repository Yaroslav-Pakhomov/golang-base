package goroutinChanWaitGroupOther

import (
	"fmt"
	"sync"
)

// region Теоретическая часть

// Горутины почти всегда используют каналы, WaitGroup или context:
// 	- каналы — для данных
// 	- sync.WaitGroup — для ожидания
// 	- context — для отмены и таймаутов

// Иначе они могут просто не успеть выполниться. Если main завершился — вся программа завершается, даже если горутины ещё работают.
// То есть:
// 	- нет «фоновых» горутин после main
// 	- main — такая же горутина, просто главная

// Как запомнить
// 			Нужно 					Используй
// 	Передать данные				chan
// 	Дождаться завершения		WaitGroup
// 	И то и другое				chan + WaitGroup

// endregion Теоретическая часть

// region Практическая часть

func GetChanWaitGroupBase() {

	// 1. Горутина + канал (chan)
	chStr := make(chan string)

	go createChan(chStr)
	msg := <-chStr

	fmt.Println(msg)

	// 2. Несколько горутин + sync.WaitGroup
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)           // говорим: +1 горутина
		go createWg(i, &wg) // запускаем
	}

	wg.Wait() // ждём, пока все вызовут Done
	fmt.Println("Горутины отработали.")

	// 3. Несколько горутин + канал (классический паттерн)
	chInt := make(chan int)

	for i := 0; i < 3; i++ {
		go workerChan(i, chInt)
	}

	fmt.Println("")
	for i := 0; i < 3; i++ {
		result := <-chInt
		fmt.Println("Квадрат числа", result)
	}

	// Задание 1. Закреп канал
	chStr2 := make(chan string)

	go writeChan(chStr2)

	resultStr := <-chStr2

	fmt.Println("")
	fmt.Println(resultStr)

	// Задание 2. Квадрат числа
	number := 5
	chInt2 := make(chan int)

	go squareChan(number, chInt2)

	resFinal := <-chInt2
	fmt.Println("")
	fmt.Printf("Квадрат числа %d равен %d.", number, resFinal)
	fmt.Println("")
	fmt.Println("")

	// Задание 3. Пять работников
	var wgWorker sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wgWorker.Add(1)
		go startWorker(i, &wgWorker)
	}

	wgWorker.Wait()
	fmt.Println("Все работники закончили.")
	fmt.Println("")

	// Задание 4. Каналы + несколько горутин. Умножение
	chTen := make(chan int)

	for i := 10; i <= 15; i++ {
		go getTenMult(i, chTen)
	}

	for i := 0; i <= 4; i++ {
		result := <-chTen
		fmt.Printf("Умножение на 10 равен %d.", result)
		fmt.Println("")
	}

	// Задание 5. Канал + WaitGroup
	// Алгоритм:
	// 	- WaitGroup ждёт производителей (writers)
	// 	- отдельная горутина закрывает канал после Wait()
	// 	- main читает через range до закрытия
	// 	- ничего не потеряется и не “успеет/не успеет”

	var wgChan sync.WaitGroup
	chInt3 := make(chan int)

	for i := 1; i <= 5; i++ {
		wgChan.Add(1)

		go getChanWaitGroup(i, chInt3, &wgChan)
	}

	// Закрывает канал после Wait()
	go closeChanWaitGroup(chInt3, &wgChan)

	fmt.Println("")
	for val := range chInt3 {
		fmt.Println("Результат квадрата числа:", val)
	}

}

// region Горутина + канал (chan)

// Канал — это труба для передачи данных между горутинами.
func createChan(ch chan string) {
	ch <- "Привет из канала"
}

// Что важно понять
// 	- make(chan string) — канал для строк
// 	- ch <- "..." — отправка в канал
// 	- <-ch — получение (и ожидание!)
// 	- main блокируется, пока данные не придут

// Поэтому Sleep тут уже не нужен — канал сам синхронизирует.

// endregion Горутина + канал (chan)

// region Несколько горутин + sync.WaitGroup

func createWg(id int, wg *sync.WaitGroup) {
	defer wg.Done() // сообщаем: горутина завершилась
	fmt.Println("Работает горутина с", id)
}

// Как это работает
// 	- wg.Add(1) — увеличиваем счётчик
// 	- wg.Done() — уменьшаем (обычно через defer)
// 	- wg.Wait() — блокируем main, пока счётчик ≠ 0

// Очень надёжный способ ожидания.

// endregion Несколько горутин + sync.WaitGroup

// region Несколько горутин + канал (классический паттерн)

func workerChan(number int, chInt chan int) {
	chInt <- number * 2
}

// Ключевая идея
// 	- main знает, сколько результатов ждать
// 	- каждый <-ch блокируется, пока данные не придут
// 	- порядок вывода не гарантирован

// endregion Несколько горутин + канал (классический паттерн)

// region Задание 1. Закреп канал

// Условие:
// 	- Запусти одну горутину
// 	- Она должна отправить в канал строку "Горутина завершилась"
// 	- main должен вывести это сообщение

func writeChan(chStr chan string) {
	chStr <- "Горутина завершилась"
}

// endregion Задание 1. Закреп канал

// region Задание 2. Квадрат числа

// Условие:
// 	- В main есть число n
// 	- Горутина считает n * n
// 	- Результат возвращается через канал и печатается

// Подсказка:
// Тип канала — chan int

func squareChan(number int, chInt chan int) {
	chInt <- number * number
}

// endregion Задание 2. Квадрат числа

// region Задание 3. Пять работников

// Условие:
// 	- Запусти 3 горутины
// 	- Каждая печатает: Worker X started
// 	- После завершения всех — main пишет: All workers done

// Требования:
// 	- Используй sync.WaitGroup
// 	- wg.Done() через defer

func startWorker(number int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Работник", number, "начал.")
}

// endregion Задание 3. Пять работников

// region Задание 4. Каналы + несколько горутин. Умножение

// Условие:
// 	- Запусти 4 горутины
// 	- Каждая получает число i
// 	- Отправляет в канал i * 10
// 	- main читает ровно 4 значения и печатает

// Важно:
// 	- Порядок вывода не важен
// 	- Sleep запрещён

func getTenMult(number int, chInt chan int) {
	chInt <- number * 10
}

// endregion Задание 4. Каналы + несколько горутин. Умножение

// region Задание 5. Канал + WaitGroup

// Условие:
// 	- Несколько горутин:
// 		* делают работу
// 		* пишут результат в канал
// 	- WaitGroup ждёт их завершения
// 	- После wg.Wait():
// 		* канал закрывается
// 		* main читает канал через for range

// Подсказка:
// Канал закрывает только тот, кто его создавал (обычно main).

func getChanWaitGroup(number int, chInt chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	chInt <- number * number
}

func closeChanWaitGroup(chInt chan int, wg *sync.WaitGroup) {
	wg.Wait()
	close(chInt)
}

// endregion Задание 5. Канал + WaitGroup

// endregion Практическая часть

// region Контрольные вопросы

// Что делает <-ch?
// Ответ:
// Получает значение из канала и блокируется, пока значение не появится.

// Можно ли читать из пустого канала?
// Ответ:
// Да, но горутина заблокируется, пока туда не отправят данные.

// Зачем нужен sync.WaitGroup?
// Ответ:
// Чтобы дождаться завершения набора горутин, не передавая данные.

// Что будет, если забыть wg.Done()?
// Ответ:
// wg.Wait() никогда не завершится → deadlock.

// Кто должен закрывать канал?
// Ответ:
// Отправитель, а не получатель.

// Что будет при чтении из закрытого канала?
// Ответ:
// 	- получим zero-value
// 	- чтение не блокируется

// Можно ли закрыть канал дважды?
// Ответ:
// Нет, будет panic.

// Когда лучше WaitGroup, а когда chan?
// Ответ:
// 	- WaitGroup — просто дождаться завершения
// 	- chan — передать данные
// 	- оба — если нужно и то, и другое

// Гарантирует ли канал порядок?
// Ответ:
// от одной горутины — да
// от нескольких — нет

// endregion Контрольные вопросы

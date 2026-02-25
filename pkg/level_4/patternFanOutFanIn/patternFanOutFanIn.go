package patternFanOutFanIn

import (
	"fmt"
	"sync"
)

// Fan-out / Fan-in (распараллеливание)

// region Теоретическая часть

// Типовые алгоритмы с горутинами
// Fan-out — разбили задачу на части
// Fan-in — собрали результаты

// Пример простой с пояснениями:
// Fan-out — одну большую задачу раскидываем на несколько горутин
// Fan-in — собираем результаты из нескольких горутин в одном месте

// Типичный сценарий:
// 👉 есть список данных →
// 👉 обрабатываем их параллельно →
// 👉 собираем результат

// Пример задачи:
// Допустим, у нас есть числа, и мы хотим параллельно возвести их в квадрат, а потом собрать результаты.

// endregion Теоретическая часть

// region Практическая часть

func GetPatternFanOutFanIn() {
	var jobs = make(chan int)
	var results = make(chan int)

	var wg sync.WaitGroup

	// --- FAN-OUT ---
	// запускаем несколько воркеров
	numWorkers := 3
	wg.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		go workers(i, jobs, results, &wg)
	}

	// Отправляем задания
	go getJobs(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	// --- FAN-IN ---
	// собираем результаты
	for result := range results {
		fmt.Println("Результат:", result)
	}
}

// workers — это Fan-out часть, читает числа из jobs, считает квадрат и пишет в results
func workers(id int, jobs <-chan int, result chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for j := range jobs {
		fmt.Println("Worker", id, "обрабатывает", j)
		result <- j * j
	}
}

// getJobs - отправляем задания
func getJobs(jobs chan<- int) {
	for i := 1; i <= 5; i++ {
		jobs <- i
	}

	close(jobs) // больше заданий не будет
}

// region Разбор примера

// 1️⃣ Fan-out — распараллеливание
// 	for i := 1; i <= numWorkers; i++ {
// 		go worker(i, jobs, results, &wg)
// 	}

// запускаем несколько горутин-воркеров
// все они читают из одного канала jobs
// Go сам балансирует, кто какое задание возьмёт

// -----------------

// 2️⃣ Рассылка задач
// 	for i := 1; i <= 5; i++ {
// 		jobs <- i
// 	}
// 	close(jobs)

// отправляем задания в канал
// close(jobs) — сигнал воркерам: «работы больше не будет»

// -----------------

// 3️⃣ Fan-in — сбор результатов
// 	for result := range results {
// 		fmt.Println("Результат:", result)
// 	}

// один потребитель
// читает из канала results
// канал закрывается только когда все воркеры завершились

// -----------------

// Визуальная схема
// 	  jobs
// 		↓
// ┌───────────┐
// │ worker 1  │
// ├───────────┤
// │ worker 2  │  ← FAN-OUT
// ├───────────┤
// │ worker 3  │
// └───────────┘
// 		↓
// 	 results
// 		↓
// 	  FAN-IN

// -----------------

// Где это реально используют:
// 	- параллельная обработка файлов;
// 	- HTTP-запросы к нескольким сервисам;
// 	- обработка очередей;
// 	- CPU-bound задачи (хеширование, парсинг, расчёты).

// endregion Разбор примера

// endregion Практическая часть

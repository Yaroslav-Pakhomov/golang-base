package selectBase

import (
	"context"
	"fmt"
	"time"
)

// Select

// region Теоретическая часть

// Что такое select

// select — как switch, но по операциям с каналами:
// 	- ждёт, какая операция (чтение/запись) сможет выполниться
// 	- выполняет одну готовую ветку
// 	- если готово несколько — выбирает случайно (fair-ish)
// 	- если ничего не готово:
// 		* без default — блокируется
// 		* с default — выполняет default сразу (не блокируется)

// -----------------

// Зачем нужен select

// Если у тебя есть несколько источников событий (несколько каналов), то без select часто получаются:
// 	- либо последовательные чтения (которые блокируют программу, если первый канал «молчит»),
// 	- либо сложные конструкции с отдельными горутинами “на каждый канал”.

// select решает задачу: подожди, пока будет готов хотя бы один из каналов, и обработай событие.

// -----------------

// Базовый синтаксис

// 	select {
// 	case v := <-ch1:
// 	   // обработка v из ch1
// 	case v := <-ch2:
// 	   // обработка v из ch2
// 	default:
// 	   // если ни один case не готов прямо сейчас (не блокируемся)
// 	}

// Важно:
// 	- select блокируется, пока не станет готов хотя бы один case, если нет default.
// 	- Если готово несколько веток, Go выбирает псевдослучайно (чтобы не было постоянного перекоса).
// 	- default превращает select в неблокирующую проверку (“polling”). С ним легко случайно сделать активное ожидание (busy loop).

// -----------------

// Чтение из закрытого канала: «нулевое значение» и ok

// Если канал закрыт, чтение:
// 	v, ok := <-ch

// даёт:
// 	- ok == false
// 	- v — нулевое значение типа (0, "", nil, …)

// Это критично в select: закрытый канал считается готовым всегда, поэтому если ты продолжишь читать закрытый канал, ты можешь получить бесконечный поток нулевых значений и «забить» цикл.

// -----------------

// Правильная обработка закрытия каналов в select

// Есть два распространённых подхода.

// Подход A: отключать канал, присваивая ему nil
// Чтение из nil-канала блокируется всегда, поэтому ветка исчезает из соревнования.

// 	for ch1 != nil || ch2 != nil {
// 	   select {
// 	   case v, ok := <-ch1:
// 	       if !ok { ch1 = nil; continue }
// 	       _ = v
// 	   case v, ok := <-ch2:
// 	       if !ok { ch2 = nil; continue }
// 	       _ = v
// 	   }
// 	}

// Подход B: отдельная логика завершения
// Например, заранее знать, сколько значений придёт, или иметь done-канал, который закрывается при завершении.

// -----------------

// done-канал для остановки горутин и отмены работы

// Паттерн:
// 	- есть канал done типа chan struct{} (часто его закрывают для broadcast-сигнала “stop”);
// 	- рабочие горутины в select слушают done.

// 	select {
// 	case <-done:
// 	   return
// 	case v := <-ch:
// 	   // работа
// 	}

// Почему struct{}? Он не занимает память для значения, а нам важен лишь факт события.

// -----------------

// Таймауты и периодические события

// Таймаут:
// 	select {
// 	case v := <-ch:
// 		_ = v
// 	case <-time.After(200 * time.Millisecond):
// 		// не дождались
// 	}

// Минус: time.After создаёт таймер каждый раз. В цикле чаще лучше time.NewTimer.

// Тикер:
// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ticker.C:
// 			// раз в секунду
// 		case <-done:
// 			return
// 		}
// 	}

// -----------------

// Типичные ошибки и как их избегать

// 1. Busy loop из-за default

// 	for {
// 		select {
// 		case v := <-ch:
// 			_ = v
// 		default:
// 			// пусто -> крутится на 100% CPU
// 		}
// 	}

// Решение: убери default, либо добавь time.Sleep, либо используй таймер/тикер.

// 2. Чтение из закрытого канала вечно
// 	Если канал закрыт и ты не отключил его (через nil) — select будет часто выбирать эту ветку.
// 3. Утечки горутин
// 	Если горутина ждёт чтения, а никто не пишет, и нет done/контекста — она может висеть вечно.

// Самые частые ошибки с select
// 	- Забыли default, думали “проверю и пойду дальше” → а оно блокируется.
// 	- Не обработали закрытие канала → бесконечно читают нули.
// 	- В цикле много time.After → лишние таймеры/нагрузка.

// endregion Теоретическая часть

// region Практическая часть

func GetSelectBase() {

	// Базовый пример
	getBaseExample()

	// Тайм-аут
	fmt.Println("")
	getTimeout()

	// Non-blocking
	fmt.Println("")
	getNonBlocking()

	// Отмена через context
	getStopBySignal()

	// Закрытие канала + select (правильная обработка)
	fmt.Println("")
	getCloseChan()

	// Работа с NewTimer
	fmt.Println("")
	getTimer()

	// Работа Fan-in
	fmt.Println("")
	getFanIn()
}

// -----------------

// getBaseExample - Базовый пример: кто первый — того и читаем
func getBaseExample() {
	chStr1 := make(chan string)
	chStr2 := make(chan string)

	go func() {
		time.Sleep(200 * time.Millisecond)
		chStr1 <- "Строка 1"
	}()
	go func() {
		time.Sleep(100 * time.Millisecond)
		chStr2 <- "Строка 2"
	}()

	select {
	case msg1 := <-chStr1:
		fmt.Println(msg1)
	case msg2 := <-chStr2:
		fmt.Println(msg2)
	}
}

// Здесь почти всегда победит chStr2, потому что он пришлёт раньше.

// -----------------

// getTimeout - тайм-аут: ждём определение время
func getTimeout() {
	chStr1 := make(chan string)

	select {
	case msg1 := <-chStr1:
		fmt.Println(msg1)
	case <-time.After(200 * time.Millisecond):
		fmt.Println("Ожидание закончилось.")
	}
}

// 💡 Если хочешь делать это в цикле — time.After может быть не лучшим (создаёт таймер каждый раз). Там лучше time.NewTimer.

// -----------------

// getNonBlocking - Non-blocking: “проверь и иди дальше”
func getNonBlocking() {
	ch2 := make(chan int)

	select {
	case msg := <-ch2:
		fmt.Println(msg)
	default:
		fmt.Println("Каналов для чтения нет.")
	}
}

// -----------------

// Отмена через context: “остановись по сигналу”

// Очень реальный паттерн: в select слушаем и работу, и ctx.Done().

// getStopBySignal - Очень реальный паттерн
func getStopBySignal() {
	jobs := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())

	go worker(ctx, jobs)

	jobs <- 1
	jobs <- 2

	cancel()
	time.Sleep(2 * time.Second)
}

// 	- jobs := make(chan int)
// 	  Канал без буфера. Это важно:
// 		* отправка jobs <- 1 блокирует getStopBySignal, пока воркер не прочитает.

// 	- ctx, cancel := context.WithCancel(...)
// 	  Создаётся контекст, который можно отменить вручную через cancel().

// 	- go worker(ctx, jobs)
// 	  Запуск воркера в отдельной горутине.

// 	- jobs <- 1, jobs <- 2
// 	  getStopBySignal отправляет 1 и 2. Поскольку канал без буфера, это синхронизируется с чтением воркером:
// 		* getStopBySignal не перейдёт к jobs <- 2, пока воркер не примет 1
// 		* и не перейдёт к cancel(), пока воркер не примет 2

// 	- cancel()
// 	  Отмена контекста → канал ctx.Done() закрывается → воркер при следующем select увидит Done и выйдет.

// 	- time.Sleep(50ms)
// 	  Это “костыль”, чтобы getStopBySignal не завершился раньше, чем воркер успеет вывести stop: ....
// 	  Без sleep программа может завершиться сразу после cancel(), и вы не увидите печать.

func worker(ctx context.Context, jobs <-chan int) {
	for {
		select {
		case job := <-jobs:
			fmt.Println("Работа", job)

		case <-ctx.Done():
			fmt.Println("Стоп", ctx.Err())
			return
		}
	}
}

// 	- for { ... } — бесконечный цикл: воркер постоянно ждёт либо работу, либо отмену.

// 	- select — ждёт готовности одного из кейсов:
// 		* case <-ctx.Done()
// 		  Сработает, когда контекст отменят (или истечёт дедлайн).
// 		  Тогда:
// 			- печатает stop: ...
// 			- делает return → горутина завершается
// 		* case job := <-jobs
// 		  Сработает, когда из канала jobs можно прочитать значение.
// 		  Печатает job: N.

// Важно: select выбирает готовый кейс. Если готовы оба — выбирает псевдослучайно.

// -----------------

// Закрытие канала + select (правильная обработка)

// Когда канал закрывают, чтение сразу готово и возвращает “нулевое значение”.
// Поэтому часто нужно читать с ok.

// getCloseChan - Закрытие канала + select
func getCloseChan() {
	ch := make(chan int, 2)
	ch <- 1
	close(ch)

	for {
		select {

		case val, ok := <-ch:
			if !ok {
				fmt.Println("Канал закрыт, выход.")
				return
			}
			fmt.Println("Получили значение", val)
		}
	}
}

// -----------------

// Работа с NewTimer. Частая “мина”: time.After в цикле → лучше NewTimer
func getTimer() {
	ch := make(chan int)

	timer := time.NewTimer(200 * time.Millisecond)
	// Останавливаем таймер, когда метод отработает
	defer timer.Stop()

	for {
		select {
		case val, ok := <-ch:
			if !ok {
				fmt.Println("канал закрыт, выходим")
				return
			}
			fmt.Println(val)
		case <-timer.C:
			fmt.Println("Время лимита вышло.")
			return
		}
	}
}

// -----------------

// Работа Fan-in (сведение нескольких каналов в один) — супер-паттерн
func getFanIn() {
	chInt1 := make(chan string)
	chInt2 := make(chan string)
	out := make(chan string)

	// fan-in
	go func() {
		defer close(out)

		for chInt1 != nil || chInt2 != nil {
			select {
			case v, ok := <-chInt1:
				if !ok {
					// Отключаем case
					chInt1 = nil
					continue
				}
				out <- v + " A"
			case v, ok := <-chInt2:
				if !ok {
					// Отключаем case
					chInt2 = nil
					continue
				}
				out <- v + " В"
			}
		}
	}()

	// Запись в первый канал
	go func() {
		chInt1 <- "Строка 1"
		time.Sleep(100 * time.Millisecond)
		chInt1 <- "Строка 2"
		// Закрывает, чтобы горутина fan-in не была бесконечной
		close(chInt1)
	}()

	// Запись во второй канал
	go func() {
		chInt2 <- "Строка 3"
		time.Sleep(100 * time.Millisecond)
		chInt2 <- "Строка 4"
		// Закрывает, чтобы горутина fan-in не была бесконечной
		close(chInt2)
	}()

	for val := range out {
		fmt.Println(val)
	}
}

// endregion Практическая часть

// region Работа с time.NewTimer

// Базовый пример: подождать один раз

// 	timer := time.NewTimer(2 * time.Second)
// 	<-timer.C // ждём пока сработает таймер, блокируется ~200s
// 	fmt.Println("Сработало через 2 секунды")

// C у time.Timer — это канал, в который таймер отправляет событие, когда истекает время.

// -----------------

// Таймер + select (часто используется)

// 	timer := time.NewTimer(2 * time.Millisecond)
// 	defer timer.Stop()

// 	select {
// 	case <-timer.C:
// 		fmt.Println("Таймаут!")
// 	}

// -----------------

// Остановить таймер (важно)

// Stop() возвращает true, если таймер был успешно остановлен до срабатывания.

// 	timer := time.NewTimer(5 * time.Second)

// 	if timer.Stop() {
// 		fmt.Println("Таймер остановлен до срабатывания")
// 	} else {
// 		// если он уже сработал — в канале может быть значение
// 		select {
// 		case <-timer.C: // вычитываем старое событие
// 		default:
// 		}
// 	}

// -----------------

// Перезапуск таймера через Reset

// Важно: перед Reset() нужно корректно остановить таймер и (иногда) очистить канал.

// 	timer := time.NewTimer(2 * time.Second)
//
// 	for i := 0; i < 3; i++ {
// 		<-timer.C
// 		fmt.Println("тик", i)
//
// 		// безопасный reset
// 		timer.Reset(1 * time.Second)
// 	}

// Если ты хочешь “таймер, который можно часто сбрасывать” (например debounce/timeout) — скажи, и я дам готовый шаблон под твой кейс.

// endregion Работа с time.NewTimer

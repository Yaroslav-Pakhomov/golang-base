package main

import (
	"golang-base/pkg/level_5/crudApiUser"
	"golang-base/pkg/level_5/httpClient"
	"time"
)

func main() {

	// region 1-ый этап

	// fmt.Println("Hello Go.")
	// fmt.Println("My name is Иван.")
	// fmt.Println("I’m starting to learn Go!")

	// // Переменные и конкатенация строк
	// greeting.Greeting()
	// // Ввод данных от пользователя
	// askname.Askname()

	// Арифметические действия
	// mathAct.Summa()
	// mathAct.Calculator()
	// mathAct.ArithmeticMean()

	// Переменные разных типов и константы
	// varConst.SetVars()
	// varConst.SetConst()

	// Оператор if / else
	// ifElse.GetEven()

	// Циклы for
	// loop.GetLoop()
	// loop.GetPracticeLoop()

	// Функция с несколькими возвращаемыми значениями
	// function.RequestFunctions()
	// Функция деления с обработкой деления на ноль и вывод ошибки
	// function.CheckFunc()

	// Оператор switch
	// switchCase.GetDay("Mon")
	// switchCase.GetDay("Thu")
	// fmt.Println("")

	// switchCase.GetGraduate(5)
	// switchCase.GetGraduate(3)
	// switchCase.GetGraduate(2)
	// fmt.Println("")

	// switchCase.GetTemperature(-2)
	// switchCase.GetTemperature(2)
	// switchCase.GetTemperature(12)
	// switchCase.GetTemperature(22)
	// fmt.Println("")

	// switchCase.CheckX(1)
	// switchCase.CheckX(2)
	// switchCase.CheckX(22)
	// fmt.Println("")

	// switchCase.GetRune('a')
	// switchCase.GetRune('b')
	// switchCase.GetRune('2')

	// Анонимная функция
	// funcAnonym.GetAnonymFuncs()

	// Функция с именованными возвращаемые значения
	// funcNamedVal.GetFuncsWithNamedVals()

	// Работа с массивами - подсчёт суммы элементов массива
	// arrElemSumma.GetArrayElementSumma()

	// Работа со срезами (Slice) - динамическими списками/массивами
	// arrSlice.WorkWithSlice()

	// Работа с картами (Map) ассоциативными массивами
	// arrMap.WorkWithMap()

	// Работа со структурами ("классами")
	// structAndMeth.GetStructs()

	// endregion 1-ый этап

	// region 2-ой этап

	// Указатели
	// pointer.WorkWithPointer()

	// Работа со структурами ("классами")
	// Работа с интерфейсами
	// structAndMethDeep.GetStructDeep()

	// Использование интерфейса
	// interfaceWork.GetAllStruct()

	// Композиция структур
	// structComposition.GetMainWork()

	// Пустой интерфейс
	// emptyInterface.GetWorkCheck()

	// Структура с JSON-тегами и Сериализация и десериализация структуры в JSON
	// structJson.GetWorkStructs()

	// endregion 2-ой этап

	// region 3-ий этап

	// Основное определение errors
	// errorsBase.GetTestErrors()

	// Пользовательская errors
	// errorsCustom.GetTestCustomErrors()

	// Логирование
	// logging.GetTestLogs()
	// logging.GetTestLogsFile()

	// endregion 3-ий этап

	// region 4-ый этап

	// Горутины база
	// goroutinBase.GetGoroutineBase()

	// Горутины с каналом, с WaitGroup, нескольких горутин
	// goroutinChanWaitGroupOther.GetChanWaitGroupBase()

	// Буферизированный канал
	// bufferChan.GetBufferChanBase()

	// Select
	// selectBase.GetSelectBase()

	// Паттерн Fan-out / Fan-in (распараллеливание)
	// patternFanOutFanIn.GetPatternFanOutFanIn()

	// Паттерн Producer / Consumer (производитель / потребитель)
	// patternProducerConsumer.GetPatternProducerConsumer()

	// Паттерн Worker Pool (пул работников)
	// patternWorkerPool.GetPatternWorkerPool()

	// endregion 4-ый этап

	// region 5-ый этап

	// Базовый HTTP-сервер
	// JSON-ответ на GET /json
	// POST-запрос с JSON-телом
	// middleware для логирования
	// httpServerBase.GetHttpServerBase()

	// CRUD API для User
	// Подключение роутера chi
	// Context в обработчиках
	// Запускаем сервер в отдельной goroutine
	go crudApiUser.GetCrudApiUser()

	// Даём серверу время стартовать
	time.Sleep(1 * time.Second)

	// HTTP-клиент
	httpClient.GetHttpClient()

	// Блокируем main goroutine, чтобы сервер продолжал работать. Иначе main() завершится все goroutine умрут.
	select {}

	// endregion 5-ый этап
}

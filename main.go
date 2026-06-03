package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"golang-base/pkg/level_5/crudApiUser"
	"golang-base/pkg/level_5/httpClient"
	"golang-base/pkg/level_6/config"
	"golang-base/pkg/level_6/database"
	"log"
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

	// region 6-ой этап

	cfg := config.LoadConfig()

	// region CRUD по database/sql
	runSQLCRUDDemo(cfg)
	// endregion CRUD по database/sql

	printBlankLines(3)

	// region CRUD по GORM
	runGormCRUDDemo(cfg)
	// endregion CRUD по GORM

	// endregion 6-ой этап

	// region 5-ый этап

	// Базовый HTTP-сервер
	// JSON-ответ на GET /json
	// POST-запрос с JSON-телом
	// middleware для логирования
	// httpServerBase.GetHttpServerBase()

	// Запускаем HTTP-клиент в отдельной goroutine
	go func() {
		// Даём серверу время стартовать
		time.Sleep(1 * time.Second)
		// HTTP-клиент
		httpClient.GetHttpClient()
	}()

	// CRUD API для User
	// Подключение роутера chi
	// Context в обработчиках
	// Загрузка файлов через POST
	// Список файлов через GET
	// Запуск сервера с учётом graceful shutdown (остановки сервера)

	// Запускаем сервер.
	// GetCrudApiUser блокирует main goroutine до Ctrl+C,
	// потому что внутри сервера работает graceful shutdown.
	// Запуск сервера происходит после подключения к БД
	crudApiUser.GetCrudApiUser()

	// endregion 5-ый этап
}

// runSQLCRUDDemo — демонстрация CRUD через database/sql.
func runSQLCRUDDemo(cfg *config.Config) {
	ctx := context.Background()

	// region Подключение к БД
	db, err := database.ConnectPostgresDb(ctx, cfg)
	if err != nil {
		log.Fatal("sql connect:", err)
	}
	defer closeDB(db)

	if err := database.CheckConnect(ctx, db); err != nil {
		log.Fatal("sql check connect:", err)
	}
	fmt.Println("DB connected (database/sql)")
	// endregion Подключение к БД

	// region Создание Табл. Поста
	if err := database.CreatePostsTable(db); err != nil {
		log.Fatal("sql create table:", err)
	}
	// endregion Создание Табл. Поста

	// region Создание записи Поста
	// err := database.CreatePost(db, "Заголовок", "Описание", 1)
	// err = database.CreatePost(db, "Заголовок 1", "Описание 1", 2)
	// err = database.CreatePost(db, "Заголовок 2", "Описание 2", 3)
	// if err != nil {
	// 	log.Fatal("sql create post:", err)
	// }
	// endregion Создание записи Поста

	// region Получение всех Постов
	posts, err := database.SelectPosts(db)
	if err != nil {
		log.Fatal("sql select posts:", err)
	}
	printJSON(posts)
	// endregion Получение всех Постов

	// region Получение Поста
	post, err := database.GetPostById(db, 1)
	if err != nil {
		log.Fatal("sql get post:", err)
	}
	printJSON(post)
	// endregion Получение Поста

	// region Обновление Поста
	if err := database.UpdatePost(db, 1, "Обновлённый заголовок", "Обновлённое описание", 10); err != nil {
		log.Fatal("sql update post:", err)
	}

	updatedPost, err := database.GetPostById(db, 1)
	if err != nil {
		log.Fatal("sql get updated post:", err)
	}
	printJSON(updatedPost)
	// endregion Обновление Поста

	// region Удаление Поста
	// if err := database.DeletePost(db, 3); err != nil {
	// 	log.Fatal("sql delete post:", err)
	// }

	postsAfterDelete, err := database.SelectPosts(db)
	if err != nil {
		log.Fatal("sql select posts after delete:", err)
	}
	printJSON(postsAfterDelete)
	// endregion Удаление Поста
}

// runGormCRUDDemo — демонстрация CRUD через GORM.
func runGormCRUDDemo(cfg *config.Config) {
	ctx := context.Background()

	// region Подключение к БД
	gormDb, err := database.ConnectPostgresGormDb(ctx, cfg)
	if err != nil {
		log.Fatal("gorm connect:", err)
	}
	sqlDB, err := gormDb.DB()
	if err != nil {
		log.Fatal("gorm get sql.DB:", err)
	}
	defer closeDB(sqlDB)

	if err := database.CheckConnectGorm(ctx, gormDb); err != nil {
		log.Fatal("gorm check connect:", err)
	}
	fmt.Println("DB connected (GORM)")
	// endregion Подключение к БД

	// region Создание Табл. Поста
	if err := database.CreatePostsTableGorm(gormDb); err != nil {
		log.Fatal("gorm create table:", err)
	}
	// endregion Создание Табл. Поста

	// region Создание записи Поста
	// if err := database.CreatePostGorm(gormDb, "title gorm", "description gorm", 5); err != nil {
	// 	log.Fatal("gorm create post:", err)
	// }
	// endregion Создание записи Поста

	// region Получение всех Постов
	postsGorm, err := database.SelectPostsGorm(gormDb)
	if err != nil {
		log.Fatal("gorm select posts:", err)
	}
	printJSON(postsGorm)
	// endregion Получение всех Постов

	// region Обновление Поста
	if err := database.UpdatePostGorm(gormDb, 7, "title gorm update", "description gorm update", 5); err != nil {
		log.Fatal("gorm update post:", err)
	}
	// endregion Обновление Поста

	// region Удаление Поста
	// if err := database.DeletePostGorm(gormDb, 9); err != nil {
	// 	log.Fatal("gorm delete post:", err)
	// }
	// endregion Удаление Поста

	// region Получение Поста
	postGorm, err := database.GetPostByIdGorm(gormDb, 7)
	if err != nil {
		log.Fatal("gorm get post:", err)
	}
	printJSON(postGorm)
	// endregion Получение Поста
}

// printJSON — сериализация и вывод любой структуры/слайса в консоль.
func printJSON(v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Println("json marshal:", err)
		return
	}
	fmt.Println(string(data))
}

func printBlankLines(n int) {
	for range n {
		fmt.Println()
	}
}

func closeDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Println("failed to close database:", err)
	}
}

package logging

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

// Логирование

// region Теоретическая часть

// Зачем нужно логирование

// Логи помогают:
// 	- понимать, что делает программа и в каком порядке;
// 	- находить ошибки и причины падений;
// 	- анализировать работу приложения на сервере (когда нет возможности «поставить дебаггер»);
// 	- фиксировать важные события (запуск, запросы, тайминги, ошибки).

// -----------------

// Пакет log в Go: основное

// Стандартный пакет log пишет текстовые сообщения и по умолчанию выводит их в stderr (часто это удобно для диагностики).

// Базовые функции:
// 	- log.Print(), log.Println(), log.Printf() — обычный вывод.
// 	- log.Fatal(), log.Fatalln(), log.Fatalf() — печатают сообщение и завершают программу через os.Exit(1) (defer НЕ выполнятся).
// 	- log.Panic(), log.Panicln(), log.Panicf() — печатают сообщение и вызывают panic (defer выполнятся, если есть recover — можно перехватить).

// -----------------

// Формат логов: флаги

// Формат задаётся через log.SetFlags(...):

// Часто используемые флаги:
// 	- log.Ldate — дата
// 	- log.Ltime — время
// 	- log.Lmicroseconds — микросекунды
// 	- log.Lshortfile — имя файла и строка (коротко)
// 	- log.Llongfile — полный путь и строка
// 	- log.LUTC — время в UTC
// 	- log.Lmsgprefix — позволяет выделять префикс как “msg prefix”

// Пример:
// 	- log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

// Префикс:
// 	- log.SetPrefix("[APP] ")

// -----------------

// Логирование в файл

// Нужно открыть файл и указать его как output:
// 	- file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
// 	- log.SetOutput(file)

// Важно:
// 	- закрывать файл через defer file.Close();
// 	- думать о том, куда писать логи: файл, stderr, stdout, или и то, и другое.

// -----------------

// Несколько логгеров

// Можно создать отдельные логгеры:
// 	- infoLogger := log.New(os.Stdout, "INFO: ", flags)
// 	- errorLogger := log.New(os.Stderr, "ERROR: ", flags)

// Это удобно, чтобы:
// 	- разделять потоки логов;
// 	- по-разному форматировать;
// 	- писать ошибки отдельно.

// endregion Теоретическая часть

// region Практическая часть

// Задание 1. “Мини-приложение с логированием”

// Сделай консольную программу, которая:
// 	1. При старте пишет лог “Приложение запущено”.
// 	2. Принимает аргументы командной строки:
// 		* --name (имя пользователя)
// 		* --age (возраст)
// 	3. Валидирует входные данные:
// 		* если имя пустое — логировать ошибку;
// 		* если возраст ≤ 0 или слишком большой (например > 120) — логировать ошибку.
// 	4. При успешной валидации печатает приветствие.
// 	5. Пишет логи в файл app.log и в консоль.

// Требования к логам
// 	- В логах должны быть: дата, время, и короткий файл/строка.
// 	- Должны быть разные уровни хотя бы логически: INFO и ERROR (через разные Logger или префиксы).
// 	- Ошибки должны писаться в stderr.

// Проверка в консоли:
// 	- go run . --name=Ира --age=25
// 	- go run . --name= --age=25
// 	- go run . --name=Ира --age=-3

func GetTestLogs() {
	// 1) Настройка флагов логирования
	flags := log.Ldate | log.Ltime | log.Lshortfile

	// 2) Открываем файл логов (append)
	f, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// Нельзя писать в файл -> пишем в stderr и выходим
		log.Fatalf("не удалось открыть файл логов: %v", err)
	}

	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Printf("ERROR: не удалось закрыть файл логов: %v", cerr)
		}
	}()

	// 3) Два канала:
	// INFO -> stdout + file
	// ERROR -> stderr + file
	infoWriter := io.MultiWriter(os.Stdout, f)
	errorWriter := io.MultiWriter(os.Stderr, f)

	infoLog := log.New(infoWriter, "INFO: ", flags)
	errLog := log.New(errorWriter, "ERROR: ", flags)

	// 4) Аргументы
	name := flag.String("name", "", "имя пользователя")
	age := flag.Int("age", 0, "возраст пользователя")
	flag.Parse()

	infoLog.Println("Приложение запущено")

	// 5) Валидация
	if *name == "" {
		errLog.Println("имя не задано (параметр --name)")
		os.Exit(1)
	}

	if *age <= 0 || *age > 120 {
		errLog.Printf("некорректный возраст: %d", *age)
		os.Exit(1)
	}

	// 6) Успешный сценарий
	infoLog.Printf("Ввод корректен: name=%s age=%d", *name, *age)
	infoLog.Printf("Привет, %s! Тебе %d лет.", *name, *age)
}

// -----------------

// Задание 2. “Логи для мини-сервиса”

// Задача: читать --input, писать непустые строки в --output, логировать счётчики.

// Сделай программу, которая:
// 	1. Читает из аргументов:
// 		* --input путь к входному файлу
// 		* --output путь к выходному файлу
// 	2. Читает входной файл построчно и пишет в выходной файл те строки, которые не пустые.
// 	3. Логирование:
// 		* INFO: старт программы, параметры, сколько строк прочитано/записано, завершение.
// 		* ERROR: если файл не открылся, если ошибка чтения/записи.
// 	4. Формат:
// 		* В логах обязателен Lshortfile.
// 		* Префиксы должны быть INFO: и ERROR: (как в модуле).
// 	5. Логи пишутся в файл app.log и в консоль.

// Критерии сдачи
// 	- Код компилируется.
// 	- При ошибках есть понятные логи.
// 	- Итоговые счетчики корректны.
// 	- app.log содержит полную историю запусков (append).

// Проверка
// go run . --input=app.log --output=app_test.log

func GetTestLogsFile() {
	flags := log.Ldate | log.Ltime | log.Lshortfile

	// Лог-файл общий для запусков
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("не удалось открыть файл логов: %v", err)
	}
	// Close() возвращает error -> обрабатываем
	defer func() {
		if cerr := logFile.Close(); cerr != nil {
			// используем обычный log, чтобы не зависеть от наших логгеров на закрытии
			log.Printf("ошибка закрытия app.log: %v", cerr)
		}
	}()

	infoWriter := io.MultiWriter(os.Stdout, logFile)
	errorWriter := io.MultiWriter(os.Stderr, logFile)

	infoLog := log.New(infoWriter, "INFO: ", flags)
	errLog := log.New(errorWriter, "ERROR: ", flags)

	// --- args ---
	inputPath := flag.String("input", "", "путь к входному файлу")
	outputPath := flag.String("output", "", "путь к выходному файлу")
	flag.Parse()

	infoLog.Println("Старт программы")
	infoLog.Printf("Параметры: input=%q output=%q", *inputPath, *outputPath)

	if *inputPath == "" || *outputPath == "" {
		errLog.Println("оба параметра обязательны: --input и --output")
		os.Exit(1)
	}

	// --- open input ---
	in, err := os.Open(*inputPath)
	if err != nil {
		errLog.Printf("не удалось открыть input файл: %v", err)
		os.Exit(1)
	}
	defer func() {
		if cerr := in.Close(); cerr != nil {
			errLog.Printf("ошибка закрытия input файла: %v", cerr)
		}
	}()

	// --- open output ---
	out, err := os.OpenFile(*outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		errLog.Printf("не удалось открыть output файл: %v", err)
		os.Exit(1)
	}
	// Важно: Flush должен выполниться ДО Close.
	// Defer выполняются в обратном порядке, поэтому:
	// 1) Ставим defer Close первым
	// 2) Ставим defer Flush вторым -> он выполнится раньше Close
	defer func() {
		if cerr := out.Close(); cerr != nil {
			errLog.Printf("ошибка закрытия output файла: %v", cerr)
		}
	}()

	writer := bufio.NewWriter(out)
	defer func() {
		if ferr := writer.Flush(); ferr != nil {
			errLog.Printf("ошибка flush writer: %v", ferr)
		}
	}()

	// --- processing ---
	scanner := bufio.NewScanner(in)

	var readLines, writtenLines int

	for scanner.Scan() {
		readLines++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		if _, err := writer.WriteString(line + "\n"); err != nil {
			errLog.Printf("ошибка записи в output: %v", err)
			os.Exit(1)
		}
		writtenLines++
	}

	if err := scanner.Err(); err != nil {
		errLog.Printf("ошибка чтения input: %v", err)
		os.Exit(1)
	}

	infoLog.Printf("Готово. Прочитано строк: %d, записано непустых: %d", readLines, writtenLines)
	infoLog.Println("Завершение программы")
}

// endregion Практическая часть

// region Контрольные вопросы — ответы

// Чем отличаются Print*, Fatal* и Panic*?
// Print* — печатает и продолжает. Fatal* — печатает и делает os.Exit(1) (defer не выполнится). Panic* — печатает и вызывает panic (defer выполнится, можно recover).

// Почему Fatal* может быть опасен при наличии defer?
// Fatal опасен, потому что делает os.Exit(1) → никакие defer не отработают (файлы не закроются, буферы могут не сброситься).

// В какой поток по умолчанию пишет пакет log?
// По умолчанию log пишет в stderr.

// Как добавить в лог дату и время? Как добавить файл и строку?
// Дата/время: log.SetFlags(log.Ldate | log.Ltime). Файл и строка: log.Lshortfile (или log.Llongfile).

// Зачем создавать разные log.Logger, если есть log.Print?
// Разные log.Logger нужны, чтобы развести потоки (stdout/stderr), префиксы/формат, разные outputs.

// Чем отличается log.SetOutput(...) от log.New(...)?
// log.SetOutput(w) меняет вывод глобального логгера. log.New(w, prefix, flags) создаёт отдельный логгер (удобно для INFO/ERROR).

// Как реализовать запись логов одновременно в консоль и в файл?
// Одновременно в консоль и файл: io.MultiWriter(os.Stdout, file) (или os.Stderr для ошибок).

// Почему полезно писать ошибки в stderr, а инфо в stdout?
// INFO в stdout, ERROR в stderr — удобно для пайпов/перенаправлений и мониторинга (ошибки отдельно от обычного вывода).

// endregion Контрольные вопросы — ответы

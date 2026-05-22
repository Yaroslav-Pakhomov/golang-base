package httpServerBase

// Базовый HTTP-сервер
// JSON-ответ на GET /json
// POST-запрос с JSON-телом
// Middleware для логирования

// Для проверки требуется:
// - Запустить "go run ." в одном терминале
// - Запуск "curl http://localhost:8080/hello" в др. терминале
// - Запуск "curl http://localhost:8080/json" в др. терминале
// - Запуск "curl -X POST http://localhost:8080/message -H "Content-Type: application/json" -d '{"name":"Ivan"}'" в др. терминале

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// loggingMiddleware - middleware (промежуточный обработчик) для логирования HTTP-запросов.
// Он перехватывает каждый запрос, выводит информацию о нём и затем передаёт запрос дальше.
// Принимает next http.Handler, next — это следующий обработчик, который должен выполниться после middleware (middleware -> handler).
// Возвращает http.Handler, потому что middleware «оборачивает» другой handler.
func loggingMiddleware(next http.Handler) http.Handler {
	// Возвращается новый handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логирование запроса: r.Method → HTTP метод, r.URL.Path → путь
		fmt.Println(r.Method, r.URL.Path)
		// Вызов следующего Handler - передай управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}

// GetHttpServerBase - запускает HTTP-сервер на порту 8080
func GetHttpServerBase() {

	// Регистрируем обработчик для маршрута "/hello"
	// При запросе:
	// curl http://localhost:8080/hello
	// будет вызвана функция helloHandler
	http.HandleFunc("/hello", helloHandler)

	// Регистрируем HTTP-обработчик для маршрута "/json".
	// При запросе:
	// curl http://localhost:8080/json
	// будет вызвана функция jsonHandler.
	http.HandleFunc("/json", jsonHandler)

	// Регистрируем HTTP-обработчик для маршрута "/json".
	// При запросе:
	// curl -X POST http://localhost:8080/message -H "Content-Type: application/json" -d '{"name":"Ivan"}'
	// будет вызвана функция jsonHandler.
	http.HandleFunc("/message", createMessageHandler)

	// Выводим сообщение о запуске сервера
	fmt.Println("server started on :8080")

	// Берём стандартный роутер
	mux := http.DefaultServeMux

	// Оборачиваем router в middleware.
	// Используется http.DefaultServeMux, потому что регистрируем маршруты через: http.HandleFunc(...).
	// Они автоматически попадают в стандартный роутер Go: http.DefaultServeMux. Именно его middleware должен «обернуть».
	loggedMux := loggingMiddleware(mux)

	// Запускаем HTTP-сервер
	// ":8080" — сервер слушает порт 8080
	// nil — используется стандартный ServeMux (маршрутизатор)
	// loggedMux с middleware. Порт можно менять нап.: 8080 -> 8081
	err := http.ListenAndServe(":8080", loggedMux)

	// Если сервер не смог запуститься — выводим ошибку
	if err != nil {
		fmt.Println("server error:", err)
	}
}

// helloHandler — обработчик HTTP-запроса
// w — ответ клиенту;
// r — запрос от клиента.
// Это фундамент всего веба в Go.

// Handler:
// - читает запрос;
// - валидирует данные;
// - выполняет логику;
// - отправляет ответ.
// Это основная архитектура веб-приложений.
func helloHandler(w http.ResponseWriter, r *http.Request) {

	// Отправляем текстовый ответ клиенту
	// w — объект ответа сервера
	// r — объект HTTP-запроса
	_, err := fmt.Fprintln(w, "Hello world, Go HTTP.")

	// Проверяем ошибку записи ответа
	if err != nil {
		return
	}
}

// Response — структура ответа сервера.
// Поле Message будет преобразовано в JSON-поле "message".
type Response struct {
	Message string `json:"message"`
}

// jsonHandler — HTTP-обработчик,
// который возвращает JSON-ответ.
func jsonHandler(w http.ResponseWriter, r *http.Request) {

	// Устанавливаем заголовок ответа.
	// Сообщаем клиенту, что сервер возвращает JSON.
	w.Header().Set("Content-Type", "application/json")

	// Создаём JSON и сразу записываем его в ответ клиенту.
	// json.NewEncoder(w) создаёт encoder,
	// который пишет данные прямо в ResponseWriter.
	err := json.NewEncoder(w).Encode(Response{
		// Данные для JSON-ответа
		Message: "Hello world, Go JSON.",
	})
	// кодировать JSON: json.NewEncoder(w).Encode(...)
	// декодировать JSON: json.NewDecoder(r.Body).Decode(...)

	// Проверяем ошибку кодирования/записи JSON
	if err != nil {
		return
	}
}

// CreateMessageRequest — структура входящего JSON-запроса.
// Ожидаем JSON вида:
//
//	{
//	  "name": "Ivan"
//	}
type CreateMessageRequest struct {
	Name string `json:"name"`
}

// Тег json:"name" — критически важен.

// createMessageHandler — HTTP-обработчик,
// который принимает JSON от клиента
// и возвращает JSON-ответ.
func createMessageHandler(w http.ResponseWriter, r *http.Request) {

	// Проверяем HTTP-метод
	if r.Method != http.MethodPost {

		// Возвращаем ошибку 405 Method Not Allowed
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Создаём переменную для хранения данных запроса
	var req CreateMessageRequest

	// Читаем JSON из тела HTTP-запроса (r.Body)
	// и записываем данные в структуру req.
	err := json.NewDecoder(r.Body).Decode(&req)

	// Если JSON некорректный —
	// возвращаем ошибку 400 Bad Request.
	if err != nil {
		// Отправляем клиенту текст ошибки
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Устанавливаем заголовок ответа.
	// Сервер будет возвращать JSON.
	w.Header().Set("Content-Type", "application/json")

	// Формируем JSON-ответ и отправляем клиенту.
	// map[string]string преобразуется в JSON-объект.
	errJson := json.NewEncoder(w).Encode(map[string]string{

		// Поле message в JSON-ответе
		"message": "Hello, " + req.Name,
	})

	// Проверяем ошибку записи JSON-ответа
	if errJson != nil {
		return
	}
}

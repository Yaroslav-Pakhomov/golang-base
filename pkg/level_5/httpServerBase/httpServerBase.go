package httpServerBase

// Базовый HTTP-сервер

// Для проверки требуется:
// - Запустить "go run ." в одном терминале
// - Запуск "curl http://localhost:8080/hello" в др. терминале
// - Запуск "curl http://localhost:8080/json" в др. терминале

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetHttpServerBase - запускает HTTP-сервер на порту 8080
func GetHttpServerBase() {

	// Регистрируем обработчик для маршрута "/hello"
	// При запросе http://localhost:8080/hello
	// будет вызвана функция helloHandler
	http.HandleFunc("/hello", helloHandler)

	// Регистрируем HTTP-обработчик для маршрута "/json".
	// При запросе http://localhost:8080/json
	// будет вызвана функция jsonHandler.
	http.HandleFunc("/json", jsonHandler)

	// Выводим сообщение о запуске сервера
	fmt.Println("server started on :8080")

	// Запускаем HTTP-сервер
	// ":8080" — сервер слушает порт 8080
	// nil — используется стандартный ServeMux (маршрутизатор)
	err := http.ListenAndServe(":8080", nil)

	// Если сервер не смог запуститься — выводим ошибку
	if err != nil {
		fmt.Println("server error:", err)
	}
}

// helloHandler — обработчик HTTP-запроса
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

	// Проверяем ошибку кодирования/записи JSON
	if err != nil {
		return
	}
}

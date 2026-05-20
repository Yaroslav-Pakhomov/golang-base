package httpServerBase

// Базовый HTTP-сервер

// Для проверки требуется:
// - Запустить "go run ." в одном терминале
// - Запуск "curl http://localhost:8080/hello" в др. терминале

import (
	"fmt"
	"net/http"
)

// GetHttpServerBase - запускает HTTP-сервер на порту 8080
func GetHttpServerBase() {

	// Регистрируем обработчик для маршрута "/hello"
	// При запросе http://localhost:8080/hello
	// будет вызвана функция helloHandler
	http.HandleFunc("/hello", helloHandler)

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

package crudApiUser

// CRUD API для User

// Проверка:
// curl -X POST http://localhost:8081/users -H "Content-Type: application/json" -d '{"name":"Alex","email":"alex@mail.com"}'
// curl http://localhost:8081/users
// curl http://localhost:8081/users/1
// curl -X PUT http://localhost:8081/users/1 -H "Content-Type: application/json" -d '{"name":"Alex Updated","email":"updated@mail.com"}'
// curl -X DELETE http://localhost:8081/users/1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = map[int]User{}
var nextID = 1

// Защита от одновременного доступа к данным.
var mu sync.RWMutex

func loggingMiddlewareCrud(next http.Handler) http.Handler {
	// Возвращается новый handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логирование запроса: r.Method → HTTP метод, r.URL.Path → путь
		fmt.Println(r.Method, r.URL.Path)
		// Вызов следующего Handler - передай управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}

func GetCrudApiUser() {
	// Подключаем роутер chi
	router := chi.NewRouter()

	// Оборачиваем router в middleware.
	router.Use(loggingMiddlewareCrud)

	router.Get("/users", getUsersHandler)
	router.Get("/users/{id}", getUserHandler)
	router.Post("/users", createUserHandler)
	router.Put("/users/{id}", updateUserHandler)
	router.Delete("/users/{id}", deleteUserHandler)

	// Запускаем HTTP-сервер, с middleware
	err := http.ListenAndServe(":8081", router)

	// Если сервер не смог запуститься — выводим ошибку
	if err != nil {
		fmt.Println("server error:", err)
	}
}

// region Методы CRUD

// getUsersHandler - получение пользователей
func getUsersHandler(w http.ResponseWriter, r *http.Request) {

	// Блокировка для чтения
	// Пока один handler держит lock — остальные ждут.
	mu.RLock()
	// Когда функция закончится — доступ откроется следующему handler.
	defer mu.RUnlock()

	result := []User{}

	for _, user := range users {
		result = append(result, user)
	}

	// Устанавливаем заголовок ответа.
	// Сообщаем клиенту, что сервер возвращает JSON.
	writeJSON(w, http.StatusOK, result)
}

// getUserHandler - получение пользователя
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getUserIDFromPath(r)
	if err != nil {
		// неправильный формат запроса → 400
		// объекта с таким ID нет → 404
		// StatusNotFound - 404, StatusBadRequest - 400
		writeError(w, http.StatusBadRequest, "Invalid user id")
		return
	}

	mu.RLock()
	defer mu.RUnlock()

	user, ok := users[id]
	if !ok {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// createUserHandler - создание пользователя
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.Name == "" || req.Email == "" {
		writeError(w, http.StatusBadRequest, "name and email required")
		return
	}

	// Блокировка для записи
	mu.Lock()
	defer mu.Unlock()

	user := User{
		ID:    nextID,
		Name:  req.Name,
		Email: req.Email,
	}
	users[nextID] = user
	nextID++

	writeJSON(w, http.StatusCreated, user)
}

// updateUserHandler - обновление пользователя
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getUserIDFromPath(r)

	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid user id")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, ok := users[id]; !ok {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}

	var req CreateUserRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.Name == "" || req.Email == "" {
		writeError(w, http.StatusBadRequest, "name and email required")
		return
	}

	user := User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}
	users[id] = user

	writeJSON(w, http.StatusOK, user)
}

// deleteUserHandler - удаление пользователя
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {

	id, err := getUserIDFromPath(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid user id")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, ok := users[id]; !ok {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}

	delete(users, id)

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "User deleted",
	})
}

// region Методы-хэлперы

// getUserIDFromPath - получение ID из запроса
func getUserIDFromPath(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")

	return strconv.Atoi(idStr)
}

// writeError - helper для отправки JSON-ошибок клиенту.
// Используется вместо http.Error, чтобы все ошибки были в одном формате JSON.
// w http.ResponseWriter - Объект ответа HTTP.
// status - HTTP статус ответа:
// - 400 Bad Request
// - 404 Not Found
// - 500 Internal Server Error
// message - текст ошибки
func writeError(w http.ResponseWriter, status int, message string) {
	// Вызываем универсальный helper writeJSON.
	// Создаём map:
	// {
	//   "error": "текст ошибки"
	// }
	// который будет автоматически преобразован в JSON.
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

// writeJSON - универсальный helper для отправки JSON-ответов.
// Может отправлять:
// - struct
// - map
// - slice
// - любые данные в формате JSON.
// data any - любые Go-данные для сериализации в JSON.
func writeJSON(w http.ResponseWriter, status int, data any) {
	// Устанавливаем HTTP header.
	// Сообщаем клиенту, что сервер отправляет JSON.
	w.Header().Set("Content-Type", "application/json")

	// Отправляем HTTP статус:
	// 200 OK
	// 201 Created
	// 400 Bad Request
	// 404 Not Found
	// и т.д.
	w.WriteHeader(status)

	// Создаём JSON encoder.
	// Encode:
	// 1. превращает Go-данные в JSON
	// 2. сразу записывает JSON в HTTP response.
	err := json.NewEncoder(w).Encode(data)

	// Проверяем ошибку сериализации JSON.
	// После WriteHeader уже нельзя менять HTTP response,
	// поэтому просто выводим ошибку в консоль сервера.
	if err != nil {
		fmt.Println("failed to encode json:", err)
	}
}

// endregion Методы-хэлперы

// endregion Методы CRUD

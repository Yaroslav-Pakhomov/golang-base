package httpServerBase

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
	"strings"
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
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userByIDHandler)

	// Выводим сообщение о запуске сервера
	// fmt.Println("server started on :8080")

	// Берём стандартный роутер
	mux := http.DefaultServeMux

	// Оборачиваем router в middleware.
	loggedMux := loggingMiddlewareCrud(mux)

	// Запускаем HTTP-сервер, с middleware
	err := http.ListenAndServe(":8081", loggedMux)

	// Если сервер не смог запуститься — выводим ошибку
	if err != nil {
		fmt.Println("server error:", err)
	}
}

// region Методы Распределения CRUD
func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsersHandler(w, r)
	case http.MethodPost:
		createUserHandler(w, r)
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

func userByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUserHandler(w, r)
	case http.MethodPut:
		updateUserHandler(w, r)
	case http.MethodDelete:
		deleteUserHandler(w, r)
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

// endregion Методы Распределения CRUD

// region Методы CRUD

// getUsersHandler - получение пользователей
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	result := []User{}

	for _, user := range users {
		result = append(result, user)
	}

	// Устанавливаем заголовок ответа.
	// Сообщаем клиенту, что сервер возвращает JSON.
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}

// getUserHandler - получение пользователя
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getUserIDFromPath(r)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusNotFound)
		return
	}

	user, ok := users[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func getUserIDFromPath(r *http.Request) (int, error) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")

	return strconv.Atoi(idStr)
}

// createUserHandler - создание пользователя
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	user := User{
		ID:    nextID,
		Name:  req.Name,
		Email: req.Email,
	}
	users[nextID] = user

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// updateUserHandler - обновление пользователя
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getUserIDFromPath(r)

	if err != nil {
		http.Error(w, "Invalid user id", http.StatusNotFound)
		return
	}

	if _, ok := users[id]; !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var req CreateUserRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	user := User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}
	users[id] = user

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// deleteUserHandler - удаление пользователя
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {

	id, err := getUserIDFromPath(r)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusNotFound)
		return
	}

	if _, ok := users[id]; !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	delete(users, id)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deleted",
	})
}

// endregion Методы CRUD

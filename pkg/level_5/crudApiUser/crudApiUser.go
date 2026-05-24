package crudApiUser

// CRUD API для User

// Проверка:
// - Запустить "go run ." в одном терминале
// curl -X POST http://localhost:8081/users -H "Content-Type: application/json" -d '{"name":"Alex","email":"alex@mail.com"}'
// curl http://localhost:8081/users
// curl http://localhost:8081/users/1
// curl -X PUT http://localhost:8081/users/1 -H "Content-Type: application/json" -d '{"name":"Alex Updated","email":"updated@mail.com"}'
// curl -X DELETE http://localhost:8081/users/1

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-base/pkg/level_5/serverCrud"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	// Проверяем/создаём папку uploads
	errDir := os.MkdirAll("./uploads", os.ModePerm)
	if errDir != nil {
		log.Fatal("cannot create uploads directory:", errDir)
	}

	// Подключаем роутер chi
	router := chi.NewRouter()

	// Оборачиваем router в middleware.
	router.Use(loggingMiddlewareCrud)

	// CRUD Users
	router.Get("/users", getUsersHandler)
	router.Get("/users/{id}", getUserHandler)
	router.Post("/users", createUserHandler)
	router.Put("/users/{id}", updateUserHandler)
	router.Delete("/users/{id}", deleteUserHandler)

	// Список файлов
	router.Get("/files", filesHandler)
	// Загрузка файла
	router.Post("/upload", uploadHandler)

	// Запуск сервера
	serverCrud.StartCrudServer(router)
}

// region Методы CRUD

// getUsersHandler - получение пользователей
func getUsersHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	// Проверяем, не был ли отменён context запроса.
	//
	// checkContext вернёт ошибку если:
	// - клиент закрыл соединение;
	// - request был отменён;
	// - истёк timeout/deadline.
	//
	// Важно:
	// Context проверяется ДО выполнения основной логики,
	// чтобы не тратить ресурсы сервера на уже отменённый запрос.
	if err := checkContext(ctx); err != nil {

		// Отправляем HTTP ошибку клиенту.
		//
		// err.Error() может вернуть:
		// - "context canceled"
		// - "context deadline exceeded"
		//
		// StatusRequestTimeout (408) сообщает,
		// что запрос был прерван или превысил timeout.
		writeError(w, http.StatusRequestTimeout, err.Error())
		return
	}

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

	ctx := r.Context()
	if err := checkContext(ctx); err != nil {
		writeError(w, http.StatusRequestTimeout, err.Error())
		return
	}

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

	ctx := r.Context()
	if err := checkContext(ctx); err != nil {
		writeError(w, http.StatusRequestTimeout, err.Error())
		return
	}

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

	ctx := r.Context()
	if err := checkContext(ctx); err != nil {
		writeError(w, http.StatusRequestTimeout, err.Error())
		return
	}

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

	ctx := r.Context()
	if err := checkContext(ctx); err != nil {
		writeError(w, http.StatusRequestTimeout, err.Error())
		return
	}

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

type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// filesHandler - возвращает список файлов из папки ./uploads.
func filesHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем содержимое папки ./uploads
	//
	// os.ReadDir возвращает:
	// - список файлов/папок
	// - ошибку
	entries, err := os.ReadDir("./uploads")

	// Если не удалось прочитать папку:
	// - папка не существует
	// - нет прав доступа
	// - ошибка файловой системы
	if err != nil {
		writeError(w, http.StatusInternalServerError, "cannot read uploads")
		return
	}

	// Создаём пустой slice для хранения имён файлов
	var files []FileInfo

	// Перебираем все элементы папки
	for _, entry := range entries {

		// Проверяем:
		// является ли элемент НЕ директорией
		//
		// То есть пропускаем вложенные папки
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			// Добавляем имя файла в slice
			files = append(files, FileInfo{
				Name: entry.Name(),
				Size: info.Size(),
			})
		}
	}

	// Возвращаем JSON:
	//
	// [
	//   {
	//     "name": "photo.png",
	//     "size": 12345
	//   },
	// ]
	writeJSON(w, http.StatusOK, files)
}

// uploadHandler - загрузка файла
func uploadHandler(w http.ResponseWriter, r *http.Request) {

	// Получаем файл из multipart/form-data
	//
	// "file" — имя поля формы:
	//
	// <input type="file" name="file">
	//
	// FormFile возвращает:
	// - file   → содержимое файла (поток io.Reader)
	// - header → метаданные файла
	// - err    → ошибка
	file, header, err := r.FormFile("file")

	// Если файл не был отправлен —
	// возвращаем HTTP 400 Bad Request
	if err != nil {
		writeError(w, http.StatusBadRequest, "file is required")
		return
	}

	// Закрываем файл после завершения функции
	defer file.Close()

	// Создаём файл на сервере:
	//
	// ./uploads/
	// └── имя_файла
	//
	// header.Filename — оригинальное имя файла
	filename := filepath.Base(header.Filename)
	dst, err := os.Create("./uploads/" + filename)

	// Если не удалось создать файл:
	// - нет папки uploads
	// - нет прав доступа
	// - ошибка диска
	if err != nil {
		writeError(w, http.StatusInternalServerError, "cannot save file")
		return
	}

	// Закрываем файл после завершения функции
	defer dst.Close()

	// Копируем содержимое uploaded файла
	// в файл на диске
	//
	// file → источник
	// dst  → куда сохраняем
	_, err = io.Copy(dst, file)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "cannot copy file")
		return
	}

	// Возвращаем JSON-ответ:
	//
	// {
	//   "filename": "photo.png"
	// }
	writeJSON(w, http.StatusCreated, map[string]string{
		"filename": filename,
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

// checkContext проверяет, не был ли отменён context запроса.
//
// Context в Go используется для:
// - отмены операций;
// - timeout/deadline;
// - передачи request-scoped данных;
// - управления жизненным циклом запроса.
//
// Context особенно важен для:
// - HTTP handlers;
// - SQL запросов;
// - внешних API;
// - goroutine;
// - gRPC;
// - Kafka/Redis clients.
//
// -------------------------------------------------------------------
// Основные типы context:
//
//  1. context.Background()
//     Базовый корневой context.
//     Обычно используется в main(), init(), tests.
//
//     ctx := context.Background()
//
//  2. context.TODO()
//     Используется как временная заглушка,
//     когда ещё не решили какой context нужен.
//
//     ctx := context.TODO()
//
//  3. context.WithCancel()
//     Позволяет вручную отменить context.
//
//     ctx, cancel := context.WithCancel(parent)
//     defer cancel()
//
//  4. context.WithTimeout()
//     Автоматически отменяет context через заданное время.
//
//     ctx, cancel := context.WithTimeout(parent, 5*time.Second)
//
//  5. context.WithDeadline()
//     Отменяет context в конкретный момент времени.
//
//     ctx, cancel := context.WithDeadline(parent, deadline)
//
//  6. context.WithValue()
//     Позволяет передавать request-scoped данные:
//     requestID, userID, traceID и т.д.
//
//     ctx := context.WithValue(parent, "requestID", "abc123")
//
// -------------------------------------------------------------------
// В HTTP сервере context обычно берут из:
//
//	r.Context()
//
// Такой context автоматически отменяется если:
// - клиент закрыл соединение;
// - истёк timeout;
// - сервер завершил request.
//
// -------------------------------------------------------------------
// Важно:
//
// checkContext НЕ должен:
// - создавать timeout;
// - писать HTTP response;
// - содержать бизнес-логику.
//
// Его задача — только проверить состояние context.
//
// -------------------------------------------------------------------
// Обычно timeout создают через middleware:
//
//	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
//	defer cancel()
//
// -------------------------------------------------------------------
// Функция возвращает:
//
// - nil                 → context активен;
// - context.Canceled    → запрос отменён;
// - context.DeadlineExceeded → timeout истёк.
//
// -------------------------------------------------------------------
// В текущем CRUD API проверка context скорее учебная,
// потому что операции с map выполняются мгновенно.
//
// Но в production context критически важен
// для работы с БД, API и конкурентными операциями.
func checkContext(ctx context.Context) error {

	select {
	case <-ctx.Done():
		// ctx.Done() — это канал, который закрывается,
		// когда context отменён.
		//
		// ctx.Err() возвращает причину отмены:
		// - context.Canceled
		// - context.DeadlineExceeded
		return ctx.Err()
	default:
		// default делает select неблокирующим.
		// Если context ещё активен, функция сразу вернёт nil.
		return nil
	}
}

// Теоретически это уместно, когда handler может выполнять долгую работу: запрос в БД, внешний API, чтение файла, ожидание goroutine. В текущем CRUD на map проверка context больше учебная, потому что операции выполняются мгновенно.

// endregion Методы-хэлперы

// endregion Методы CRUD

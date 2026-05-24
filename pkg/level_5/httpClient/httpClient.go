package httpClient

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func GetHttpClient() {

	// Создаём HTTP-клиент
	client := http.Client{
		// Максимальное время ожидания запроса
		// Если сервер не ответит за 5 секунды —
		// будет ошибка timeout
		Timeout: time.Second * 5,
	}

	// Отправляем GET-запрос на локальный сервер
	resp, errResp := client.Get("http://localhost:8081/users")

	// Проверяем ошибку запроса
	// Например:
	// - сервер недоступен
	// - timeout
	// - ошибка сети
	if errResp != nil {
		log.Fatal(errResp)
	}

	// Закрываем Body после завершения функции
	// Это обязательно, чтобы не было утечки ресурсов
	defer resp.Body.Close()

	// Читаем тело ответа полностью
	body, errBody := io.ReadAll(resp.Body)

	// Проверяем ошибку чтения ответа
	if errBody != nil {
		log.Fatal(errBody)
	}

	// Выводим ответ сервера в консоль
	fmt.Println(string(body))

}

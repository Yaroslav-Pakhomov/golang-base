package serverCrud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// StartCrudServer - запуск сервера
func StartCrudServer(router http.Handler) {

	// Создаём HTTP-сервер
	srv := &http.Server{

		// Адрес и порт сервера
		Addr: ":8081",

		// Router / Handler сервера
		Handler: router,
	}

	// Запускаем сервер в отдельной goroutine
	go func() {

		// Сообщение о старте сервера
		fmt.Println("server started on :8081")

		// Запускаем HTTP-сервер
		//
		// ListenAndServe блокирует goroutine
		// и начинает принимать HTTP-запросы
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {

			// ErrServerClosed возникает при нормальном Shutdown
			// и НЕ считается ошибкой.
			//
			// Поэтому проверяем:
			// err != http.ErrServerClosed
			log.Fatal(err)
		}
	}()

	// Канал сигналов ОС
	quit := make(chan os.Signal, 1)

	// Подписка на Ctrl+C и SIGTERM
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Ждём сигнал остановки
	<-quit

	fmt.Println("shutting down server...")

	// Context с timeout для graceful shutdown
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	// Graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("server stopped")
}

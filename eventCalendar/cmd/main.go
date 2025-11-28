package main

import (
	"database/sql"
	"eventCalendar/internal/config"
	"eventCalendar/internal/server"
	"eventCalendar/internal/storage"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Полноценный HTTP-сервис (описан в задаче 18 уровня 2).
// В качестве финального проекта реализуйте сервис, самостоятельно продумав детали (структуру хранения данных, формат запросов/ответов) в соответствии с заданием.
// Требования
// Фоновый воркер через канал: при создании события с напоминанием — кладём задачу в канал, воркер должен следить за временем и слать напоминания
// Чистка событий: отдельная горутина, каждые X минут должна переносить в архив старые события
// Асинхронный логгер: HTTP-хендлеры не должны писать в stdout напрямую, а должны класть записи в канал, который обрабатывает отдельная горутина
// Результат: директория с кодом сервера, инструкцией по запуску (README), примерами запросов и тестами внутри.

func main() {
	var dsn string
	conf := config.NewConfig("config/dev.yml")
	fmt.Println(conf)

	flag.StringVar(&dsn, "dsn", conf.DB.URL, "Postgres DSN")
	flag.Parse()
	if err := runMigrations(dsn, conf.Migrations.FilePath); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	log.Println("migrations done")

	db := storage.New(&conf.DB)
	server := server.New(conf, db)

	go func() {
		err := server.HttpServer.ListenAndServe()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done

}

func runMigrations(dsn, filePath string) error {
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		filePath,
		"postgres",
		driver,
	)

	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	log.Println("migrations applied")
	return nil
}

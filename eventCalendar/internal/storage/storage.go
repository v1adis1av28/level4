package storage

import (
	"context"
	"eventCalendar/internal/config"
	"eventCalendar/internal/models"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/quay/zlog"
)

type Storage struct {
	DB *pgx.Conn
}

func New(confDb *config.DBConfig) *Storage {
	conn, err := pgx.Connect(context.Background(), confDb.URL)
	if err != nil {
		panic(err)
	}
	return &Storage{DB: conn}
}

func (s *Storage) CreateEvent(event *models.Event) error {
	query := "INSERT INTO EVENTS(NAME,DESCRIPTION,STATUS,EVENT_DATE) VALUES($1,$2,$3,$4);"
	_, err := s.DB.Exec(context.Background(), query, event.Name, event.Description, event.Status, event.Date)
	if err != nil {
		return fmt.Errorf("error on operating exec insert into events, event: %v", event)
	}

	zlog.Info(context.Background()).Msgf("Event succesfully created. Item : %v", event)
	return nil
}

func (s *Storage) UpdateEvent(req *models.EventModificationRequest) error {
	isEventExist, err := s.IsEventExist(req.ID)
	if err != nil && !isEventExist {
		return fmt.Errorf("event with id: %v not found", req.ID)
	}

	date, _ := time.Parse("2006-01-02", req.NEW_DATE) // Необрабатываем ошибку при парсе так-как на этот шаг программа зайдет только полсе валидации в хендлере
	query := "UPDATE EVENTS SET "

	if len(req.NEW_DATE) != 0 && len(req.NEW_NAME) != 0 {
		query += "EVENT_DATE = $1, NAME = $2 WHERE ID = $3;"
		_, err := s.DB.Exec(context.Background(), query, date, req.NEW_NAME, req.ID)
		if err != nil {
			return err
		}
	} else if len(req.NEW_DATE) != 0 && len(req.NEW_NAME) < 1 {
		query += "EVENT_DATE = $1 WHERE ID = $2;"
		_, err := s.DB.Exec(context.Background(), query, date, req.ID)
		if err != nil {
			return err
		}
	} else {
		query += "NAME = $1 WHERE ID = $2;"
		_, err := s.DB.Exec(context.Background(), query, req.NEW_NAME, req.ID)
		if err != nil {
			return err
		}
	}

	zlog.Info(context.Background()).Msgf("Succesfully updated event:  %v", req.ID)
	return nil

}

func (s *Storage) IsEventExist(id int) (bool, error) {
	var exist bool
	fmt.Println(exist)

	row := s.DB.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM EVENTS WHERE ID = $1);", id)
	err := row.Scan(&exist)
	if err != nil {
		return false, err
	}

	fmt.Println(exist)
	return exist, nil
}

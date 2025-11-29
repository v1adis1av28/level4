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
	query := "INSERT INTO EVENTS(NAME,DESCRIPTION,STATUS,EVENT_DATE,HAVE_NOTIFICATION) VALUES($1,$2,$3,$4,$5);"
	_, err := s.DB.Exec(context.Background(), query, event.Name, event.Description, event.Status, event.Date, event.HaveNotification)
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

func (s *Storage) DeleteEvent(id int) error {
	exists, err := s.IsEventExist(id)
	if err != nil {
		return fmt.Errorf("error checking event existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("event with id %d does not exist", id)
	}

	query := "DELETE FROM EVENTS WHERE ID = $1;"
	_, err = s.DB.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("error executing delete query: %w", err)
	}
	zlog.Info(context.Background()).Msgf("event with ID %d deleted successfully", id)
	return nil
}

func (s *Storage) GetEventsForDay(dateStr string) ([]models.Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing date string '%s': %w", dateStr, err)
	}

	startOfDay := date
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)

	query := `
        SELECT ID, NAME, DESCRIPTION, EVENT_DATE, STATUS, HAVE_NOTIFICATION
        FROM EVENTS
        WHERE EVENT_DATE >= $1 AND EVENT_DATE <= $2
        ORDER BY EVENT_DATE ASC;
    `

	rows, err := s.DB.Query(context.Background(), query, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("error executing query for events on day %s: %w", dateStr, err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Date, &event.Status, &event.HaveNotification)
		if err != nil {
			return nil, fmt.Errorf("error scanning event row: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return events, nil
}

func (s *Storage) GetEventsForWeek(dateStr string) ([]models.Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing date string '%s': %w", dateStr, err)
	}

	daysToSubtract := (int(date.Weekday()) - 1 + 7) % 7
	startOfWeek := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()).AddDate(0, 0, -daysToSubtract)
	endOfWeek := startOfWeek.AddDate(0, 0, 7).Add(-time.Nanosecond)

	query := `
        SELECT ID, NAME, DESCRIPTION, EVENT_DATE, STATUS, HAVE_NOTIFICATION
        FROM EVENTS
        WHERE EVENT_DATE >= $1 AND EVENT_DATE <= $2
        ORDER BY EVENT_DATE ASC;
    `

	rows, err := s.DB.Query(context.Background(), query, startOfWeek, endOfWeek)
	if err != nil {
		return nil, fmt.Errorf("error executing query for events in week starting %s: %w", startOfWeek.Format("2006-01-02"), err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Date, &event.Status, &event.HaveNotification)
		if err != nil {
			return nil, fmt.Errorf("error scanning event row: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return events, nil
}

func (s *Storage) GetEventsForMonth(dateStr string) ([]models.Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing date string '%s': %w", dateStr, err)
	}

	year, month, _ := date.Date()
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, date.Location())
	endOfMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, date.Location()).Add(-time.Nanosecond) // 23:59:59.999...

	query := `
        SELECT ID, NAME, DESCRIPTION, EVENT_DATE, STATUS, HAVE_NOTIFICATION
        FROM EVENTS
        WHERE EVENT_DATE >= $1 AND EVENT_DATE <= $2
        ORDER BY EVENT_DATE ASC;
    `

	rows, err := s.DB.Query(context.Background(), query, startOfMonth, endOfMonth)
	if err != nil {
		return nil, fmt.Errorf("error executing query for events in month %s: %w", startOfMonth.Format("2006-01"), err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Date, &event.Status, &event.HaveNotification)
		if err != nil {
			return nil, fmt.Errorf("error scanning event row: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return events, nil
}

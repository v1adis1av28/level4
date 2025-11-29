package storage

import (
	"context"
	"database/sql"
	"eventCalendar/internal/config"
	"eventCalendar/internal/models"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	tcm "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDatabase(t *testing.T) (string, func()) {
	provider, err := tc.NewDockerProvider()
	require.NoError(t, err)
	defer provider.Close()

	pgContainer, err := tcm.Run(
		context.Background(),
		"docker.io/postgres:15-alpine",
		tcm.WithDatabase("testdb"),
		tcm.WithUsername("testuser"),
		tcm.WithPassword("testpass"),
		tc.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)

	connectionString, err := pgContainer.ConnectionString(context.Background(), "sslmode=disable")
	require.NoError(t, err)

	tempDB, err := sql.Open("postgres", connectionString)
	require.NoError(t, err)
	defer tempDB.Close()

	driver, err := postgres.WithInstance(tempDB, &postgres.Config{})
	require.NoError(t, err)
	migrationPath := "file://./migrations"
	if _, err := os.Stat("./migrations"); os.IsNotExist(err) {
		migrationPath = "file://../../migrations"
		if _, err := os.Stat("../../migrations"); os.IsNotExist(err) {
			t.Fatalf("Migrations directory not found at ./migrations or ../../migrations. Current working directory: %s", os.Getenv("PWD"))
		}
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres",
		driver,
	)
	require.NoError(t, err)

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}
	m.Close()

	cleanup := func() {
		if err := pgContainer.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate postgres container: %s", err)
		}
	}

	return connectionString, cleanup
}

func TestStorageIntegration(t *testing.T) {
	connectionString, cleanup := setupTestDatabase(t)
	defer cleanup()

	storageInstance := New(&config.DBConfig{URL: connectionString})
	defer storageInstance.DB.Close(context.Background())

	t.Run("CreateAndReadEvent", func(t *testing.T) {
		now := time.Now()
		futureDate := now.Add(24 * time.Hour).Format("2006-01-02")
		eventToCreate := &models.Event{
			Name:             "Integration Test Event",
			Description:      "This is an event created during integration test",
			Date:             futureDate,
			Status:           models.StatusScheduled,
			HaveNotification: true,
		}

		err := storageInstance.CreateEvent(eventToCreate)
		require.NoError(t, err)
		assert.NotZero(t, eventToCreate.ID)

		retrievedEvent, err := storageInstance.GetEventByID(eventToCreate.ID)
		require.NoError(t, err)

		assert.Equal(t, eventToCreate.Name, retrievedEvent.Name)
		assert.Equal(t, eventToCreate.Description, retrievedEvent.Description)
		assert.Equal(t, eventToCreate.Date, retrievedEvent.Date)
		assert.Equal(t, eventToCreate.Status, retrievedEvent.Status)
		assert.Equal(t, eventToCreate.HaveNotification, retrievedEvent.HaveNotification)
		assert.Equal(t, eventToCreate.ID, retrievedEvent.ID)
	})

	t.Run("UpdateEvent", func(t *testing.T) {
		now := time.Now()
		futureDate := now.Add(24 * time.Hour).Format("2006-01-02")
		originalEvent := &models.Event{
			Name:             "Original Name",
			Description:      "Original Description",
			Date:             futureDate,
			Status:           models.StatusScheduled,
			HaveNotification: false,
		}
		err := storageInstance.CreateEvent(originalEvent)
		require.NoError(t, err)

		newName := "Updated Name"
		newDate := now.Add(48 * time.Hour).Format("2006-01-02")
		req := &models.EventModificationRequest{
			ID:       originalEvent.ID,
			NEW_NAME: newName,
			NEW_DATE: newDate,
		}
		err = storageInstance.UpdateEvent(req)
		require.NoError(t, err)

		updatedEvent, err := storageInstance.GetEventByID(originalEvent.ID)
		require.NoError(t, err)

		assert.Equal(t, newName, updatedEvent.Name)
		assert.Equal(t, newDate, updatedEvent.Date)
		assert.Equal(t, originalEvent.Description, updatedEvent.Description)
		assert.Equal(t, originalEvent.Status, updatedEvent.Status)
		assert.Equal(t, originalEvent.HaveNotification, updatedEvent.HaveNotification)
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		now := time.Now()
		futureDate := now.Add(24 * time.Hour).Format("2006-01-02")
		eventToDelete := &models.Event{
			Name:             "To Be Deleted",
			Description:      "This event will be deleted",
			Date:             futureDate,
			Status:           models.StatusScheduled,
			HaveNotification: false,
		}
		err := storageInstance.CreateEvent(eventToDelete)
		require.NoError(t, err)

		err = storageInstance.DeleteEvent(eventToDelete.ID)
		require.NoError(t, err)

		_, err = storageInstance.GetEventByID(eventToDelete.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("ArchiveOldEvents", func(t *testing.T) {
		now := time.Now()
		oldDate := now.AddDate(0, -2, 0).Format("2006-01-02")
		futureDate := now.Add(24 * time.Hour).Format("2006-01-02")

		oldEvent1 := &models.Event{
			Name:             "Old Completed Event 1",
			Description:      "Completed long ago",
			Date:             oldDate,
			Status:           models.StatusCompleted,
			HaveNotification: false,
		}
		oldEvent2 := &models.Event{
			Name:             "Old Cancelled Event 2",
			Description:      "Cancelled long ago",
			Date:             oldDate,
			Status:           models.StatusCancelled,
			HaveNotification: false,
		}
		newEvent := &models.Event{
			Name:             "New Scheduled Event 3",
			Description:      "Scheduled for future",
			Date:             futureDate,
			Status:           models.StatusScheduled,
			HaveNotification: false,
		}

		err := storageInstance.CreateEvent(oldEvent1)
		require.NoError(t, err)
		err = storageInstance.CreateEvent(oldEvent2)
		require.NoError(t, err)
		err = storageInstance.CreateEvent(newEvent)
		require.NoError(t, err)

		cutoffTime := now.AddDate(0, -1, 0)
		err = storageInstance.ArchiveOldEvents(context.Background(), cutoffTime)
		require.NoError(t, err)

		_, err = storageInstance.GetEventByID(oldEvent1.ID)
		assert.Error(t, err)
		_, err = storageInstance.GetEventByID(oldEvent2.ID)
		assert.Error(t, err)

		retrievedNewEvent, err := storageInstance.GetEventByID(newEvent.ID)
		require.NoError(t, err)
		assert.Equal(t, newEvent.Name, retrievedNewEvent.Name)

		var count int
		err = storageInstance.DB.QueryRow(context.Background(), "SELECT COUNT(*) FROM ARCHIVE WHERE EVENT_ID = $1 OR EVENT_ID = $2", oldEvent1.ID, oldEvent2.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 2, count)

		err = storageInstance.DB.QueryRow(context.Background(), "SELECT COUNT(*) FROM ARCHIVE WHERE EVENT_ID = $1", newEvent.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

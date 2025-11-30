package tests

import (
	"database/sql"
	"demo/internal/models"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost:5433/advertisements_test?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	clearDB()
	db.Close()
	os.Exit(code)
}

func clearDB() {
	_, err := db.Exec("TRUNCATE TABLE ORDERS CASCADE")
	if err != nil {
		panic(err)
	}
	fmt.Println("db cleared")
}

func TestInsertOrder(t *testing.T) {
	_, err := db.Exec(`INSERT INTO orders(order_uid, track_number, entry) VALUES('test123', 'track123', 'web')`)
	if err != nil {
		t.Fatal(err)
	}

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM orders WHERE order_uid='test123'`).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("ожидали 1 запись, получили %d", count)
	}
}

func TestInsertDuplicateOrder(t *testing.T) {
	_, err := db.Exec(`INSERT INTO orders(order_uid, track_number, entry) VALUES('test123', 'track123', 'web')`)
	if err == nil {
		t.Fatal("Cannot insert duplicate order")
	}
}

func TestGetOrderById(t *testing.T) {
	var order models.Order
	err := db.QueryRow(`SELECT order_uid from orders where order_uid = 'test123';`).Scan(&order.OrderUID)
	if err != nil {
		t.Fatal(err)
	}
	if order.OrderUID != "test123" {
		t.Fatal("wrong value of order")
	}
}

func TestGetNotExistingOrder(t *testing.T) {
	var order models.Order
	err := db.QueryRow(`SELECT order_uid from orders where order_uid = 'test1231';`).Scan(&order.OrderUID)
	if err == nil {
		t.Fatal("order not exist")
	}
}

func TestInvalidInputOnOrder(t *testing.T) {
	invalidOrderUUID := 32
	invalidTrackNumber := "122"
	_, err := db.Exec(`INSERT INTO orders(order_uid, track_number, entry) VALUES($1, $2, 'web')`, invalidOrderUUID, invalidTrackNumber)
	if err != nil {
		t.Fatal("invalid parameters on inserting order")
	}
}

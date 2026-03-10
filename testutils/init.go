package testutils

import (
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	gConn  *gorm.DB
	mockdb sqlmock.Sqlmock
	once   sync.Once
)

func initDB() {
	c, m, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	conn, err := gorm.Open(postgres.New(postgres.Config{
		Conn: c,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	gConn = conn
	mockdb = m
}

func Mock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	once.Do(func() {
		initDB()
	})
	return gConn, mockdb
}

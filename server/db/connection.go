package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"server/errs"
)

type Connector interface {
	Open(host, port, user, password, database string) error
	Close() error
	GetConnection() (*sql.DB, error)
}

type connector struct {
	connection *sql.DB
	mu         sync.Mutex
}

func (c *connector) Open(host, port, user, password, database string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connection != nil {
		return errors.New(errs.DatabaseConnectionAlreadyOpenError)
	}

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}

	c.connection = db
	return nil
}

func (c *connector) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connection == nil {
		return errors.New(errs.DatabaseConnectionNotOpenError)
	}

	err := c.connection.Close()
	c.connection = nil
	return err
}

func (c *connector) GetConnection() (*sql.DB, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connection == nil {
		return nil, errors.New(errs.DatabaseConnectionNotEstablishedError)
	}

	return c.connection, nil
}

var HandlerConnector Connector = &connector{}

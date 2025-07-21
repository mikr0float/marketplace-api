package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gocraft/dbr/v2"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	conn *dbr.Connection
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	log.Printf("Connecting to PostgreSQL DB")
	conn, err := dbr.Open("postgres", dsn, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = conn.PingContext(ctx)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connection to PostgreSQL DB is successfull")
	return &PostgresStorage{conn: conn}, nil
}

func (s *PostgresStorage) NewSession(ctx context.Context) (*dbr.Session, error) {
	session := s.conn.NewSession(nil)
	err := session.PingContext(ctx)
	if err != nil {
		log.Printf("connection verification failed: %v", err)
		return nil, fmt.Errorf("connection verification failed, %w", err)
	}
	return session, nil
}

func (s *PostgresStorage) Close() error {
	return s.conn.Close()
}

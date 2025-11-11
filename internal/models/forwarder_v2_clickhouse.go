package models

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/google/uuid"
)

type ForwarderClickhouseV2 struct {
	Type     string `json:"type" enum:"clickhouse" required:"true"`
	Url      string `json:"url" required:"true"`
	Database string `json:"db" required:"true"`
	Table    string `json:"table" required:"true" pattern:"^[a-zA-Z_][a-zA-Z0-9_]*$" minLength:"1" maxLength:"64"`
	Username string `json:"user" required:"true"`
	Password string `json:"pass" required:"true"`
	UseTls   bool   `json:"tls" required:"true"`
}

var createDbQuery = `
CREATE TABLE IF NOT EXISTS %s (
	id         UUID                 NOT NULL PRIMARY KEY,
	timestamp  DateTime64(3, 'UTC') NOT NULL,
	fields     Map(String, String)  NOT NULL,
) ENGINE = MergeTree
`

var insertLogQuery = `
INSERT INTO %s
VALUES (?, ?, ?)
`

func (f *ForwarderClickhouseV2) call(ctx context.Context, record *LogRecord) error {
	var tlscfg *tls.Config
	if f.UseTls {
		tlscfg = &tls.Config{}
	} else {
		tlscfg = nil
	}

	var conn, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{f.Url},
		Auth: clickhouse.Auth{
			Database: f.Database,
			Username: f.Username,
			Password: f.Password,
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "TODO", Version: "TODO"},
			},
		},
		TLS: tlscfg,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize connection: %w", err)
	}
	defer conn.Close()

	if err := conn.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping server: %w", err)
	}

	query := fmt.Sprintf(createDbQuery, f.Table)
	if err := conn.Exec(ctx, query); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	pk, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to generate uuid: %w", err)
	}

	query = fmt.Sprintf(insertLogQuery, f.Table)
	if err := conn.Exec(ctx, query, pk, record.Timestamp, record.Fields); err != nil {
		return fmt.Errorf("failed to insert row: %w", err)
	}

	return nil
}

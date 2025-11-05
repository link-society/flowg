package models

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ForwarderClickhouseV2 struct {
	Type     string `json:"type" enum:"clickhouse" required:"true"`
	Url      string `json:"url" required:"true"`
	Database string `json:"db" required:"true"`
	Table    string `json:"table" required:"true"`
	Username string `json:"user" required:"true"`
	Password string `json:"pass" required:"true"`
}

var createDbQuery = `
CREATE TABLE IF NOT EXISTS ? (
	id         uuid                 NOT NULL PRIMARY KEY,
	timestamp  DateTime64(3, 'UTC') NOT NULL,
	fields     Map(String, String)  NOT NULL,
) ENGINE = MergeTree
`

var insertLogQuery = `
INSERT INTO ?
VALUES (toUUID(rand64()), ?, ?)
`

func (f *ForwarderClickhouseV2) call(ctx context.Context, record *LogRecord) error {
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
		TLS: &tls.Config{
			InsecureSkipVerify: true, // TODO
		},
	})
	if err != nil {
		return fmt.Errorf("failed to initialize connection: %w", err)
	}
	defer conn.Close()

	if err := conn.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping server: %w", err)
	}

	if err := conn.Exec(ctx, createDbQuery, f.Table); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	if err := conn.Exec(ctx, insertLogQuery, f.Table,
		record.Timestamp, record); err != nil {
		return fmt.Errorf("failed to insert row: %w", err)
	}

	return nil
}

package models

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/uuid"
	"link-society.com/flowg/internal/app"
)

type forwarderStateClickhouseV2 struct {
	conn        driver.Conn
	createQuery string
	insertQuery string
}

type ForwarderClickhouseV2 struct {
	Type     string `json:"type" enum:"clickhouse" required:"true"`
	Url      string `json:"url" required:"true"`
	Database string `json:"db" required:"true"`
	Table    string `json:"table" required:"true" pattern:"^[a-zA-Z_][a-zA-Z0-9_]*$" minLength:"1" maxLength:"64"`
	Username string `json:"user" required:"true"`
	Password string `json:"pass" required:"true"`
	UseTls   bool   `json:"tls" required:"true"`

	state *forwarderStateClickhouseV2
}

const createDbQuery = `
CREATE TABLE IF NOT EXISTS %s (
	id         UUID                 NOT NULL PRIMARY KEY,
	timestamp  DateTime64(3, 'UTC') NOT NULL,
	fields     Map(String, String)  NOT NULL,
) ENGINE = MergeTree
`

const insertLogQuery = `
INSERT INTO %s
VALUES (?, ?, ?)
`

func (f *ForwarderClickhouseV2) init(ctx context.Context) error {
	if f.state != nil {
		return fmt.Errorf("clickhouse state has already been initialized")
	}

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
				{Name: "FlowG", Version: app.FLOWG_VERSION},
			},
		},
		TLS: tlscfg,
	})
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to initialize connection: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close()
		return fmt.Errorf("failed to ping server: %w", err)
	}

	f.state = &forwarderStateClickhouseV2{
		conn:        conn,
		createQuery: fmt.Sprintf(createDbQuery, f.Table),
		insertQuery: fmt.Sprintf(insertLogQuery, f.Table),
	}

	return nil
}

func (f *ForwarderClickhouseV2) close(ctx context.Context) error {
	if f.state == nil || f.state.conn == nil {
		return fmt.Errorf("clickhouse forwarder hasn't been initialized")
	}

	return f.state.conn.Close()
}

func (f *ForwarderClickhouseV2) call(ctx context.Context, record *LogRecord) error {
	if f.state == nil || f.state.conn == nil {
		return fmt.Errorf("clickhouse state hasn't been properly initialized")
	}

	if err := f.state.conn.Exec(ctx, f.state.createQuery); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	pk, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to generate uuid: %w", err)
	}

	if err := f.state.conn.Exec(ctx, f.state.insertQuery, pk, record.Timestamp,
		record.Fields); err != nil {
		return fmt.Errorf("failed to insert row: %w", err)
	}

	return nil
}

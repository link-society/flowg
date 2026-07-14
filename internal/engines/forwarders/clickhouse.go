package forwarders

import (
	"context"
	"fmt"

	"crypto/tls"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/google/uuid"

	"link-society.com/flowg/internal/app"
	"link-society.com/flowg/internal/models"
)

type clickhouseRuntime struct {
	config *models.ForwarderClickhouseV2

	conn        driver.Conn
	insertQuery string
}

var _ Runtime = (*clickhouseRuntime)(nil)

const clickhouseCreateDbQuery = `
CREATE TABLE IF NOT EXISTS %s (
	id         UUID                 NOT NULL PRIMARY KEY,
	timestamp  DateTime64(3, 'UTC') NOT NULL,
	fields     Map(String, String)  NOT NULL,
) ENGINE = MergeTree
`

const clickhouseInsertLogQuery = `
INSERT INTO %s
VALUES (?, ?, ?)
`

func (rt *clickhouseRuntime) Init(ctx context.Context) error {
	var tlscfg *tls.Config
	if rt.config.UseTls {
		tlscfg = &tls.Config{}
	} else {
		tlscfg = nil
	}

	var conn, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{rt.config.Address},
		Auth: clickhouse.Auth{
			Database: rt.config.Database,
			Username: rt.config.Username,
			Password: rt.config.Password,
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

	createQuery := fmt.Sprintf(clickhouseCreateDbQuery, rt.config.Table)
	if err := conn.Exec(ctx, createQuery); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	rt.conn = conn
	rt.insertQuery = fmt.Sprintf(clickhouseInsertLogQuery, rt.config.Table)

	return nil
}

func (rt *clickhouseRuntime) Close(ctx context.Context) error {
	if rt.conn == nil {
		return fmt.Errorf("clickhouse forwarder hasn't been initialized")
	}

	return rt.conn.Close()
}

func (rt *clickhouseRuntime) Call(ctx context.Context, record *models.LogRecord) error {
	if rt.conn == nil {
		return fmt.Errorf("clickhouse state hasn't been properly initialized")
	}

	pk, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to generate uuid: %w", err)
	}

	if err := rt.conn.Exec(ctx, rt.insertQuery, pk, record.Timestamp,
		record.Fields); err != nil {
		return fmt.Errorf("failed to insert row: %w", err)
	}

	return nil
}

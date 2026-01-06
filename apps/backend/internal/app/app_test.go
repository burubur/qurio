package app

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"qurio/apps/backend/internal/config"
)

func TestNew_Success(t *testing.T) {
	// Arrange
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockVec := &MockVectorStore{}
	mockPub := &MockTaskPublisher{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := &config.Config{}

	// Act
	app, err := New(cfg, db, mockVec, mockPub, logger)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, app)
	assert.NotNil(t, app.Handler)
	assert.NotNil(t, app.SourceService)
	assert.NotNil(t, app.ResultConsumer)
}

type FakeDB struct{}

func (f *FakeDB) PingContext(ctx context.Context) error { return nil }
func (f *FakeDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row { return nil }
func (f *FakeDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) { return nil, nil }
func (f *FakeDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) { return nil, nil }
func (f *FakeDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) { return nil, nil }

func TestNew_PanicsOnInvalidDB(t *testing.T) {
	// Arrange
	mockVec := &MockVectorStore{}
	mockPub := &MockTaskPublisher{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := &config.Config{}
	
	fakeDB := &FakeDB{}

	// Act & Assert
	assert.Panics(t, func() {
		_, _ = New(cfg, fakeDB, mockVec, mockPub, logger)
	})
}
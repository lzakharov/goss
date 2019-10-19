package storage

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/lzakharov/goss/internal/domain"
	"go.uber.org/zap"
)

var (
	userColumns = []string{"id", "username", "role"}
)

func TestNewAdapter(t *testing.T) {
	logger := zap.NewExample()
	db := new(sqlx.DB)

	expected := &adapter{
		logger: logger,
		db:     db,
	}

	actual := NewAdapter(logger, db)

	require.Equal(t, expected, actual)
}

func TestAdapter_IsAlive(t *testing.T) {
	logger := zap.NewExample()
	db, _, err := sqlmock.New()
	require.NoError(t, err)

	adapter := &adapter{
		logger: logger,
		db:     sqlx.NewDb(db, "postgres"),
	}

	t.Run("alive", func(t *testing.T) {
		require.True(t, adapter.IsAlive())
	})
}

func TestAdapter_GetUser(t *testing.T) {
	logger := zap.NewExample()

	adapter := &adapter{
		logger: logger,
	}

	t.Run("existing user", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		adapter.db = sqlx.NewDb(db, "postgres")

		mock.ExpectQuery(`^SELECT (.+) FROM "user" WHERE id = (.+)$`).
			WillReturnRows(sqlmock.NewRows(userColumns).FromCSVString("1,alice,client"))

		expected := &domain.User{
			ID:       1,
			Username: "alice",
			Role:     "client",
		}

		actual, err := adapter.GetUser(1)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("nonexistent user", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		adapter.db = sqlx.NewDb(db, "postgres")

		mock.ExpectQuery(`^SELECT (.+) FROM "user" WHERE id = (.+)$`).
			WillReturnError(sql.ErrNoRows)

		_, err = adapter.GetUser(1)
		require.Error(t, err)
	})
}

func TestAdapter_GetUserByCredentials(t *testing.T) {
	logger := zap.NewExample()

	adapter := &adapter{
		logger: logger,
	}

	t.Run("valid credentials", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		adapter.db = sqlx.NewDb(db, "postgres")

		mock.ExpectQuery(`^SELECT (.+) FROM "user" WHERE username = (.+) AND password = (.+)$`).
			WillReturnRows(sqlmock.NewRows(userColumns).FromCSVString("1,alice,client"))

		credentials := &domain.Credentials{
			Username: "alice",
			Password: "password",
		}

		expected := &domain.User{
			ID:       1,
			Username: "alice",
			Role:     "client",
		}

		actual, err := adapter.GetUserByCredentials(credentials)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		adapter.db = sqlx.NewDb(db, "postgres")

		mock.ExpectQuery(`^SELECT (.+) FROM "user" WHERE username = (.+) AND password = (.+)$`).
			WillReturnError(sql.ErrNoRows)

		credentials := &domain.Credentials{
			Username: "alice",
			Password: "password",
		}

		_, err = adapter.GetUserByCredentials(credentials)
		require.Error(t, err)
	})
}

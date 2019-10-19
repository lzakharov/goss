package storage

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lzakharov/goss/internal/domain"
	"go.uber.org/zap"

	_ "github.com/lib/pq" // postgres driver
)

//NewAdapter creates a new storage adapter.
func NewAdapter(logger *zap.Logger, db *sqlx.DB) domain.Storage {
	adapter := &adapter{
		logger: logger,
		db:     db,
	}

	return adapter
}

type adapter struct {
	logger *zap.Logger
	db     *sqlx.DB
}

//IsAlive returns true if the adapter can ping it's database.
func (a *adapter) IsAlive() bool {
	return a.db.Ping() == nil
}

//GetUser gets user by id.
func (a *adapter) GetUser(userID int64) (*domain.User, error) {
	user := new(domain.User)

	if err := a.db.QueryRowx(
		getUserQuery,
		userID,
	).StructScan(user); err != nil {
		a.logger.Error("Error getting user by id!",
			zap.Int64("userID", userID),
			zap.Error(err))

		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, domain.ErrInternalStorage
	}

	return user, nil
}

//GetUserByCredentials gets user by the credentials.
func (a *adapter) GetUserByCredentials(credentials *domain.Credentials) (*domain.User, error) {
	user := new(domain.User)

	if err := a.db.QueryRowx(
		getUserByCredentialsQuery,
		credentials.Username,
		credentials.Password,
	).StructScan(user); err != nil {
		a.logger.Error("Error getting a user by the credentials!",
			zap.String("username", credentials.Username),
			zap.Error(err))

		if err == sql.ErrNoRows {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, domain.ErrInternalStorage
	}

	return user, nil
}

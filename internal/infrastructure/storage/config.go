package storage

//Config contains a storage adapter configuration.
type Config struct {
	DB *DBConfig `validate:"required"`
}

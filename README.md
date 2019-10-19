# GoSS

Go Simple Service.

## Usage

Check the [Makefile](./Makefile) for usage information.

## Configuration

| Environment Variable              | Description                                                    | Example                                                             |
|-----------------------------------|----------------------------------------------------------------|---------------------------------------------------------------------|
| ENVIRONMENT                       | Application environment (`DEV` or `PROD`)                      | DEV                                                                 |
| APP_STORAGE_DB_DRIVERNAME         | Database driver                                                | postgres                                                            |
| APP_STORAGE_DB_DSN                | Database DSN                                                   | postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable |
| APP_STORAGE_DB_MAXOPENCONNS       | Maximum number of open connections to the database             | 10                                                                  |
| APP_STORAGE_DB_MAXIDLECONNS       | Maximum number of connections in the idle connection pool      | 10                                                                  |
| APP_STORAGE_DB_CONNMAXLIFETIME    | Maximum amount of time a connection may be reused              | 5m                                                                  |
| APP_STORAGE_DB_MIGRATIONS_DIALECT | Database dialect                                               | postgres                                                            |
| APP_STORAGE_DB_MIGRATIONS_DIR     | Migrations directory                                           | migrations/dev/postgres                                             |
| APP_SECURITY_KEYPREFIX            | Key prefix for data storage                                    | auth                                                                |
| APP_SECURITY_SECRET               | Encryption key                                                 | secret                                                              |
| APP_SECURITY_ACCESSTOKENLIFETIME  | Access token lifetime                                          | 24h                                                                 |
| APP_SECURITY_REFRESHTOKENLIFETIME | Refresh token lifetime                                         | 720h                                                                |
| APP_SECURITY_REDISCLIENT_ADDR     | Redis address                                                  | redis:6379                                                          |
| APP_HTTP_ADDRESS                  | HTTP-server adapter                                            | :8080                                                               |
| APP_HTTP_READTIMEOUT              | Amount of time allowed to read the full request including body | 30s                                                                 |

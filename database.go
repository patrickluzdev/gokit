package gokit

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	sqlDB  *sql.DB
	gormDB *gorm.DB
	config Config
	driver string
}

func NewConnection(config Config) *Connection {
	driver := config.GetWithDefault("DB_DRIVER", "postgres")
	dsn := buildDSN(config, driver)

	sqlDB, err := sql.Open(driver, dsn)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	if err := sqlDB.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to ping database: %v", err))
	}

	var gormDB *gorm.DB
	switch driver {
	case "postgres":
		gormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		panic(fmt.Sprintf("Unsupported database driver: %s", driver))
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to connect GORM: %v", err))
	}

	return &Connection{
		sqlDB:  sqlDB,
		gormDB: gormDB,
		config: config,
		driver: driver,
	}
}

func (c *Connection) SQL() *sql.DB {
	return c.sqlDB
}

func (c *Connection) GORM() *gorm.DB {
	return c.gormDB
}

func buildDSN(config Config, driver string) string {
	switch driver {
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.Get("DB_USERNAME"),
			config.Get("DB_PASSWORD"),
			config.GetWithDefault("DB_HOST", "localhost"),
			config.GetWithDefault("DB_PORT", "5432"),
			config.Get("DB_DATABASE"),
		)
	default:
		return ""
	}
}

type DB struct {
	conn *Connection
}

func NewDB(config Config) Database {
	conn := NewConnection(config)
	return &DB{conn: conn}
}

func (d *DB) Connection() *sql.DB {
	return d.conn.SQL()
}

func (d *DB) Table(name string) QueryBuilder {
	return nil
}

func (d *DB) Model(model any) QueryBuilder {
	return NewGormQueryBuilder(d.conn.GORM(), model)
}

func (d *DB) Raw(sql string, args ...any) QueryBuilder {
	return nil
}

func (d *DB) Begin() (Transaction, error) {
	sqlTx, err := d.conn.SQL().Begin()
	if err != nil {
		return nil, err
	}

	gormTx := d.conn.GORM().Begin()
	if gormTx.Error != nil {
		sqlTx.Rollback()
		return nil, gormTx.Error
	}

	return NewTransaction(sqlTx, gormTx), nil
}

func (d *DB) Transaction(fn func(tx Transaction) error) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (d *DB) Migrate() error {
	return nil
}

func (d *DB) Seed() error {
	return nil
}

func (d *DB) Exec(sql string, args ...any) error {
	_, err := d.conn.SQL().Exec(sql, args...)
	return err
}

type Tx struct {
	sqlTx  *sql.Tx
	gormTx *gorm.DB
}

func NewTransaction(sqlTx *sql.Tx, gormTx *gorm.DB) Transaction {
	return &Tx{
		sqlTx:  sqlTx,
		gormTx: gormTx,
	}
}

func (t *Tx) Table(name string) QueryBuilder {
	return nil
}

func (t *Tx) Model(model any) QueryBuilder {
	return NewGormQueryBuilder(t.gormTx, model)
}

func (t *Tx) Commit() error {
	if err := t.gormTx.Commit().Error; err != nil {
		return err
	}
	return t.sqlTx.Commit()
}

func (t *Tx) Rollback() error {
	t.gormTx.Rollback()
	return t.sqlTx.Rollback()
}

func (t *Tx) SavePoint(name string) error {
	return t.gormTx.SavePoint(name).Error
}

func (t *Tx) RollbackTo(name string) error {
	return t.gormTx.RollbackTo(name).Error
}

func (t *Tx) Select(columns ...string) QueryBuilder {
	return t.Table("")
}

func (t *Tx) Where(query any, args ...any) QueryBuilder {
	return t.Table("")
}

func (t *Tx) WhereIn(column string, values []any) QueryBuilder {
	return t.Table("")
}

func (t *Tx) Join(query string, args ...any) QueryBuilder {
	return t.Table("")
}

func (t *Tx) OrderBy(column string, direction ...string) QueryBuilder {
	return t.Table("")
}

func (t *Tx) GroupBy(columns ...string) QueryBuilder {
	return t.Table("")
}

func (t *Tx) Having(query any, args ...any) QueryBuilder {
	return t.Table("")
}

func (t *Tx) Limit(limit int) QueryBuilder {
	return t.Table("")
}

func (t *Tx) Offset(offset int) QueryBuilder {
	return t.Table("")
}

func (t *Tx) Preload(query string, args ...any) QueryBuilder {
	return t.Table("")
}

func (t *Tx) Joins(query string, args ...any) QueryBuilder {
	return t.Table("")
}

func (t *Tx) Find(dest any) error {
	return nil
}

func (t *Tx) First(dest any) error {
	return nil
}

func (t *Tx) Create(value any) error {
	return nil
}

func (t *Tx) Update(values any) error {
	return nil
}

func (t *Tx) Delete() error {
	return nil
}

func (t *Tx) Count() (int64, error) {
	return 0, nil
}

func (t *Tx) Begin() (QueryBuilder, error) {
	return t, nil
}

func (t *Tx) ToSQL() (string, []any, error) {
	return "", nil, nil
}

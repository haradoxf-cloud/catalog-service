package rcpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/haradoxf-cloud/catalog-service/migration"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	"github.com/haradoxf-cloud/catalog-service/internal/app/config/section"
)

type (
	Client struct {
		_bunDB
		rawBunDB *bun.DB

		cfg section.RepositoryPostgres
	}

	_bunDB = bun.IDB
)

func (c *Client) GetRawBunDB() *bun.DB {
	return c.rawBunDB
}

func NewClient(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {
	var u url.URL
	u.Scheme = "postgres"
	u.Host = cfg.Address
	u.User = url.UserPassword(cfg.Username, cfg.Password)
	u.Path = cfg.Name
	args := make(url.Values)
	args.Set("sslmode", "disable")
	u.RawQuery = args.Encode()
	dsn := u.String()

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	sqlDB.SetMaxOpenConns(10)

	bunDB := bun.NewDB(sqlDB, pgdialect.New())

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("база не отвечает: %w", err)
	}

	log.Println("Успешно подключились к PSQL!")

	return &Client{_bunDB: bunDB, rawBunDB: bunDB, cfg: cfg}, nil
}

func (c *Client) Migrate(ctx context.Context) (oldVer, newVer int64, err error) {
	// 1. Загружаем миграции из встроенных файлов
	migrations := migrate.NewMigrations()
	if err = migrations.Discover(migration.Postgres); err != nil {
		return 0, 0, fmt.Errorf("failed to discover migrations: %w", err)
	}

	// 2. Создаем мигратор с опциями
	opts := []migrate.MigratorOption{
		migrate.WithTableName(c.cfg.MigrationTable),
		migrate.WithLocksTableName(c.cfg.MigrationTable + "_lock"),
		migrate.WithMarkAppliedOnSuccess(true),
	}
	m := migrate.NewMigrator(c.rawBunDB, migrations, opts...)

	// 3. Инициализируем таблицу миграций (если её нет)
	if err = m.Init(ctx); err != nil {
		return 0, 0, fmt.Errorf("failed to discover migrations: %w", err)
	}

	// 4. Получаем старую версию (ДО применения)
	applied, err := m.AppliedMigrations(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to discover migrations: %w", err)
	}
	for _, mg := range applied {
		v, _ := strconv.ParseInt(mg.Name, 10, 64)
		if v > oldVer {
			oldVer = v
		}
	}

	// 5. Применяем новые миграции
	group, err := m.Migrate(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to discover migrations: %w", err)
	}

	// 6. Вычисляем новую версию
	newVer = oldVer
	for _, mg := range group.Migrations {
		v, _ := strconv.ParseInt(mg.Name, 10, 64)
		if v > newVer {
			newVer = v
		}
	}

	return oldVer, newVer, nil
}

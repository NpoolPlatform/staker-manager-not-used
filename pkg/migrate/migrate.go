package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"entgo.io/ent/dialect"

	entsql "entgo.io/ent/dialect/sql"

	constant "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"

	archivementent "github.com/NpoolPlatform/archivement-manager/pkg/db/ent"
	archivementconst "github.com/NpoolPlatform/archivement-manager/pkg/message/const"

	billingent "github.com/NpoolPlatform/cloud-hashing-billing/pkg/db/ent"
	billingconst "github.com/NpoolPlatform/cloud-hashing-billing/pkg/message/const"

	ledgerent "github.com/NpoolPlatform/ledger-manager/pkg/db/ent"
	ledgerconst "github.com/NpoolPlatform/ledger-manager/pkg/message/const"

	orderent "github.com/NpoolPlatform/cloud-hashing-order/pkg/db/ent"
	orderconst "github.com/NpoolPlatform/cloud-hashing-order/pkg/message/const"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	_ "github.com/go-sql-driver/mysql" // nolint
)

const (
	keyUsername = "username"
	keyPassword = "password"
	keyDBName   = "database_name"
)

func dsn(hostname string) (string, error) {
	username := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyUsername)
	password := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyPassword)
	dbname := config.GetStringValueWithNameSpace(hostname, keyDBName)

	svc, err := config.PeekService(constant.MysqlServiceName)
	if err != nil {
		logger.Sugar().Warnw("dsb", "error", err)
		return "", err
	}

	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&interpolateParams=true",
		username, password,
		svc.Address,
		svc.Port,
		dbname,
	), nil
}

func open(hostname string) (conn *sql.DB, err error) {
	hdsn, err := dsn(hostname)
	if err != nil {
		return nil, err
	}

	conn, err = sql.Open("mysql", hdsn)
	if err != nil {
		return nil, err
	}

	// https://github.com/go-sql-driver/mysql
	// See "Important settings" section.
	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)

	return conn, nil
}

func _migrate(
	ctx context.Context,
	order *orderent.Client,
	billing *billingent.Client,
	archivement *archivementent.Client,
	ledger *ledgerent.Client,
) error {
	// Migrate payments to ledger details and general
	// Migrate commission to ledger detail and general
	return nil
}

func migrate(ctx context.Context, order, billing, archivement, ledger *sql.DB) error {
	return _migrate(
		ctx,
		orderent.NewClient(
			orderent.Driver(
				entsql.OpenDB(dialect.MySQL, order),
			),
		),
		billingent.NewClient(
			billingent.Driver(
				entsql.OpenDB(dialect.MySQL, billing),
			),
		),
		archivementent.NewClient(
			archivementent.Driver(
				entsql.OpenDB(dialect.MySQL, archivement),
			),
		),
		ledgerent.NewClient(
			ledgerent.Driver(
				entsql.OpenDB(dialect.MySQL, ledger),
			),
		),
	)
}

func Migrate(ctx context.Context) (err error) {
	logger.Sugar().Infow("Migrate", "Start", "...")
	defer func() {
		logger.Sugar().Infow("Migrate", "Done", "...", "error", err)
	}()

	// Prepare mysql instance for order / billing / ledger
	order, err := open(orderconst.ServiceName)
	if err != nil {
		return err
	}

	billing, err := open(billingconst.ServiceName)
	if err != nil {
		return err
	}

	archivement, err := open(archivementconst.ServiceName)
	if err != nil {
		return err
	}

	ledger, err := open(ledgerconst.ServiceName)
	if err != nil {
		return err
	}

	return migrate(ctx, order, billing, archivement, ledger)
}

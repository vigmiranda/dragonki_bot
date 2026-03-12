package postgres

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/vigmiranda/dragonki/internal/bootstrap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not defined")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("error connecting to database:", err)
	}

	AttachObservability(db)
	return db
}

func AttachObservability(db *gorm.DB) {
	if db == nil || db.Config == nil {
		return
	}

	baseLogger := db.Config.Logger
	if baseLogger == nil {
		baseLogger = logger.Default
	}
	if _, alreadyWrapped := baseLogger.(*observabilityLogger); alreadyWrapped {
		return
	}

	db.Config.Logger = &observabilityLogger{base: baseLogger}
}

type observabilityLogger struct {
	base logger.Interface
}

func (l *observabilityLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &observabilityLogger{base: l.unwrap().LogMode(level)}
}

func (l *observabilityLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.unwrap().Info(ctx, msg, data...)
}

func (l *observabilityLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.unwrap().Warn(ctx, msg, data...)
}

func (l *observabilityLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.unwrap().Error(ctx, msg, data...)
}

func (l *observabilityLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql := ""
	if fc != nil {
		sql, _ = fc()
	}

	bootstrap.IncAppMetric("db.query.total")
	bootstrap.AddAppMetric("db.query.duration_ms_total", elapsed.Milliseconds())

	switch detectSQLOperation(sql) {
	case "SELECT":
		bootstrap.IncAppMetric("db.query.select_total")
	case "INSERT":
		bootstrap.IncAppMetric("db.query.insert_total")
	case "UPDATE":
		bootstrap.IncAppMetric("db.query.update_total")
	case "DELETE":
		bootstrap.IncAppMetric("db.query.delete_total")
	default:
		bootstrap.IncAppMetric("db.query.other_total")
	}

	if err != nil {
		bootstrap.IncAppMetric("db.query.error_total")
	}
	if elapsed >= 200*time.Millisecond {
		bootstrap.IncAppMetric("db.query.slow_total")
	}

	l.unwrap().Trace(ctx, begin, fc, err)
}

func (l *observabilityLogger) unwrap() logger.Interface {
	if l != nil && l.base != nil {
		return l.base
	}
	return logger.Default
}

func detectSQLOperation(sql string) string {
	fields := strings.Fields(strings.TrimSpace(sql))
	if len(fields) == 0 {
		return "OTHER"
	}
	return strings.ToUpper(fields[0])
}

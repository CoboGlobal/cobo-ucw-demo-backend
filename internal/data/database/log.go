package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type GormLoggerConfig struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogLevel                  log.Level
}

var _ logger.Interface = (*gormLogger)(nil)

type SlowQueryFunc func(ctx context.Context, elapsed time.Duration, sql string, rowsAffected int64)

type gormLogger struct {
	GormLoggerConfig
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string

	slowQueryHandler SlowQueryFunc

	logger *log.Helper
}

func GormLogger(config GormLoggerConfig, logger *log.Helper) *gormLogger {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	return &gormLogger{
		GormLoggerConfig: config,
		infoStr:          infoStr,
		warnStr:          warnStr,
		errStr:           errStr,
		traceStr:         traceStr,
		traceWarnStr:     traceWarnStr,
		traceErrStr:      traceErrStr,
		logger:           logger,
	}
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	switch level {
	case logger.Silent:
		newlogger.LogLevel = log.LevelFatal
	case logger.Error:
		newlogger.LogLevel = log.LevelError
	case logger.Warn:
		newlogger.LogLevel = log.LevelWarn
	case logger.Info:
		newlogger.LogLevel = log.LevelInfo
	}
	return &newlogger
}

func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel <= log.LevelInfo {
		l.getLogger().Infof(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel <= log.LevelWarn {
		l.getLogger().Warnf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel <= log.LevelError {
		l.getLogger().Errorf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

//nolint:cyclop
func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel >= log.LevelFatal {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel <= log.LevelError && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.getLogger().Errorf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.getLogger().Errorf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel <= log.LevelWarn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.getLogger().Warnf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.getLogger().Warnf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
		if l.slowQueryHandler != nil {
			l.slowQueryHandler(ctx, elapsed, sql, rows)
		}
	case l.LogLevel <= log.LevelInfo:
		sql, rows := fc()
		if rows == -1 {
			l.getLogger().Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.getLogger().Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func (l *gormLogger) getLogger() *log.Helper {
	return l.logger
}

func (l *gormLogger) SetSlowQueryHandler(handler SlowQueryFunc) {
	l.slowQueryHandler = handler
}

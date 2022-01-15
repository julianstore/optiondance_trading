package config

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

var Log *zap.Logger

func InitLogger() (err error) {
	z := Cfg.Zap
	writeSyncer := getLogWriter(z.Filename, z.MaxSize, z.MaxBackups, z.MaxAge)
	stdoutSyncer := zapcore.AddSync(os.Stdout)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(z.Level))
	if err != nil {
		return
	}
	var core zapcore.Core
	if z.ConsoleOut {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, l),
			zapcore.NewCore(encoder, stdoutSyncer, l),
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}
	Log = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(Log) // Replace the global logger instance in the zap package, and then only need to use zap.L() call in other packages
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger receives the default log of the gin framework
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start)
		var b strings.Builder
		b.WriteString("[GIN]")
		b.WriteString(path)
		Log.Info(b.String(),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recovers the panic that may appear in the project, and uses zap to record related logs
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					Log.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					Log.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					Log.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// GormLogger struct
type GormLogger struct {
	logger.LogLevel
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = logger.Info
	return &newlogger
}

func (l *GormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	switch i[0] {
	case "sql":
		zap.L().Debug(
			"sql",
			zap.String("module", "gorm"),
			zap.String("type", "sql"),
			zap.Any("src", i[1]),
			zap.Any("duration", i[2]),
			zap.Any("sql", i[3]),
			zap.Any("values", i[4]),
			zap.Any("rows_returned", i[5]),
		)
	case "log":
		zap.L().Debug("log", zap.Any("gorm", i[2]))
	}
}

func (l *GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	switch i[0] {
	case "sql":
		zap.L().Debug(
			"sql",
			zap.String("module", "gorm"),
			zap.String("type", "sql"),
			zap.Any("src", i[1]),
			zap.Any("duration", i[2]),
			zap.Any("sql", i[3]),
			zap.Any("values", i[4]),
			zap.Any("rows_returned", i[5]),
		)
	case "log":
		zap.L().Debug("log", zap.Any("gorm", i[2]))
	}
}

func (l *GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	switch i[0] {
	case "sql":
		zap.L().Debug(
			"sql",
			zap.String("module", "gorm"),
			zap.String("type", "sql"),
			zap.Any("src", i[1]),
			zap.Any("duration", i[2]),
			zap.Any("sql", i[3]),
			zap.Any("values", i[4]),
			zap.Any("rows_returned", i[5]),
		)
	case "log":
		zap.L().Debug("log", zap.Any("gorm", i[2]))
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if err != nil {
		sql, rows := fc()
		if err.Error() == "record not found" {
			zap.L().Info(err.Error(), zap.String("sql", sql), zap.Int64("rows", rows))
		} else {
			zap.L().Error(err.Error(), zap.String("sql", sql), zap.Int64("rows", rows))
			//SentryCaptureMessage(err.Error(), sql)
		}
		return
	}
}

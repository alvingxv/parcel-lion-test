package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger

	once sync.Once
)

const folder = "logs"

func Initialize() error {
	var err error
	once.Do(func() {

		dir, _ := os.Getwd()
		logDir := dir + "/" + folder

		// Ensure log directory exists
		err = os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			return
		}

		debug := false

		debugStr := os.Getenv("DEBUG")
		if debugStr == "" || debugStr == "true" {
			debug = true
		}

		level := zap.InfoLevel

		if debug {
			level = zap.DebugLevel
		}
		// Create a custom log writer with daily rotation
		logWriter := newDailyLogWriter(logDir)

		// Create encoder config
		encoderConfig := getEncoderConfig()

		// Create a core that writes to both file and console
		core := zapcore.NewTee(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				logWriter,
				level,
			),
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoderConfig),
				logWriter,
				level,
			),
		)

		// Create the logger
		Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
		Sugar = Logger.Sugar()
	})

	return err
}

// dailyLogWriter implements a custom WriteSyncer that creates a new log file each day
type dailyLogWriter struct {
	logDir      string
	currentFile *os.File
	currentDate string
}

// newDailyLogWriter creates a new daily log writer
func newDailyLogWriter(logDir string) *dailyLogWriter {
	return &dailyLogWriter{
		logDir: logDir,
	}
}

// Write implements the WriteSyncer interface
func (d *dailyLogWriter) Write(p []byte) (n int, err error) {
	today := time.Now().Format("2006-01-02")

	// Create new file if date changed or no file exists
	if d.currentFile == nil || d.currentDate != today {
		// Close existing file if open
		if d.currentFile != nil {
			d.currentFile.Close()
		}

		// Create new log file
		filename := filepath.Join(d.logDir, fmt.Sprintf("%s.log", today))
		d.currentFile, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return 0, err
		}
		d.currentDate = today
	}

	// Write to file
	return d.currentFile.Write(p)
}

// Sync implements the WriteSyncer interface
func (d *dailyLogWriter) Sync() error {
	if d.currentFile != nil {
		return d.currentFile.Sync()
	}
	return nil
}

// getLogLevel converts string to zapcore.Level
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

// getEncoderConfig creates a custom encoder configuration
func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// Close flushes any buffered log entries and closes the current log file
func Close() {
	if Logger != nil {
		Logger.Sync()
	}
}

func LogRequest(c *fiber.Ctx, timestamp time.Time) {
	var requestBody []byte
	if c.Body() != nil {
		requestBody = c.Body()
	}

	reqHeaders := make(map[string]string)
	for k, v := range c.GetReqHeaders() {
		reqHeaders[k] = strings.Join(v, ",")
	}

	reqLog := Sugar.With(
		zap.Time("timestamp", timestamp),
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.String("remote_addr", c.IP()),
		zap.ByteString("request_body", requestBody),
		zap.Any("request_header", reqHeaders),
	)

	reqLog.Info("request")

}

func LogResponse(responseStatus int, responseBody string, responseTime time.Duration) {

	respLog := Sugar.With(
		zap.Int("status", responseStatus),
		zap.String("response_body", responseBody),
		zap.Duration("duration", responseTime),
	)

	respLog.Info("response")

}

func LogDebug(msg string) {
	Sugar.With(
		zap.String("msg:", msg),
	).Debug("debug")
}

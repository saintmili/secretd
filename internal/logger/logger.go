package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/saintmili/secretd/internal/config"
)

type Logger struct {
	enabled bool
	file    *os.File
}

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func New(cfg config.Logging) (*Logger, error) {
	if !cfg.Enabled {
		return &Logger{enabled: false}, nil
	}

	path := expandHome(cfg.File)

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, err
	}

	l := &Logger{
		enabled: true,
		file:    f,
	}

	l.Info("Logger initialized")
	return l, nil
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

func (l *Logger) log(level Level, msg string) {
	ts := time.Now().Format(time.RFC3339)

	var lvl string
	switch level {
	case DEBUG:
		lvl = "DEBUG"
	case INFO:
		lvl = "INFO"
	case WARN:
		lvl = "WARN"
	case ERROR:
		lvl = "ERROR"
	}

	fmt.Fprintf(l.file, "%s [%s] %s\n", ts, lvl, msg)
}

func (l *Logger) Debug(msg string) { l.log(DEBUG, msg) }
func (l *Logger) Info(msg string)  { l.log(INFO, msg) }
func (l *Logger) Warn(msg string)  { l.log(WARN, msg) }
func (l *Logger) Error(msg string) { l.log(ERROR, msg) }

func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

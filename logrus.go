package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/juju/errors"
	"github.com/sirupsen/logrus"
)

func ltolr(level Level) logrus.Level {
	switch level {
	case DebugLevel:
		return logrus.DebugLevel
	case InfoLevel:
		return logrus.InfoLevel
	case WarnLevel:
		return logrus.WarnLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case FatalLevel:
		return logrus.FatalLevel
	case PanicLevel:
		return logrus.PanicLevel
	}
	panic(fmt.Sprintf("Can't map level %d to logrus level", level))
}

func newLogrus(w io.Writer, lvl Level) Logger {
	l := logrus.New()
	l.SetOutput(w)
	l.SetLevel(ltolr(lvl))
	return &lLog{l}
}

type lLog struct {
	writer logrus.FieldLogger
}

// WithField returns a new Logger that always logs the specified field
func (l *lLog) WithField(key, value string) Logger {
	writer := l.writer.WithField(key, value)
	return &lLog{writer: writer}
}

// Level creates a new Entry with the specified Level
func (l *lLog) Level(lvl Level) Entry {
	switch lvl {
	case DebugLevel:
		return l.Debug()
	case InfoLevel:
		return l.Info()
	case WarnLevel:
		return l.Warn()
	case ErrorLevel:
		return l.Error()
	case FatalLevel:
		return l.Fatal()
	case PanicLevel:
		return l.Panic()
	default:
		return l.Info()
	}
}

// Debug creates a new Entry with level Debug
func (l *lLog) Debug() Entry {
	return &lEntry{logrus.DebugLevel, l.writer.WithField("time", time.Now())}
}

// Info creates a new Entry with level Info
func (l *lLog) Info() Entry {
	return &lEntry{logrus.InfoLevel, l.writer.WithField("time", time.Now())}
}

// Warn creates a new Entry with level Warn
func (l *lLog) Warn() Entry {
	return &lEntry{logrus.WarnLevel, l.writer.WithField("time", time.Now())}
}

// Error creates a new Entry with level Error
func (l *lLog) Error() Entry {
	return &lEntry{logrus.ErrorLevel, l.writer.WithField("time", time.Now())}
}

// Fatal creates a new Entry with level Fatal. Executing a log at fatal level exits the application with exit code 1.
func (l *lLog) Fatal() Entry {
	return &lEntry{logrus.FatalLevel, l.writer.WithField("time", time.Now())}
}

// Panic creates a new Entry with level Panic. Executing a log at panic level will call panic().
func (l *lLog) Panic() Entry {
	return &lEntry{logrus.PanicLevel, l.writer.WithField("time", time.Now())}
}

type lEntry struct {
	level logrus.Level
	entry *logrus.Entry
}

// Flush writes the entry as a single log statement. Optionally, a message can be added which will
// be included in the final log entry
func (l *lEntry) Flush(msg string) {
	l.entry.Logln(l.level, msg)
	if l.level == logrus.FatalLevel {
		os.Exit(1)
	}
}

// AddFields adds a range of fields to the log statement
func (l *lEntry) AddFields(fs map[string]interface{}) Entry {
	l.entry = l.entry.WithFields(fs)
	return l
}

// AddErr adds an error to the log statement. The error will have the key "err". An error stack will be included
// under the key "err_stack"
func (l *lEntry) AddErr(err error) Entry {
	msg := err.Error()
	st := errors.ErrorStack(err)
	l.entry = l.entry.WithField("err", msg)
	l.entry = l.entry.WithField("err_stack", st)
	return l
}

// AddError adds an error to the log statement. An error stack will be included under the key "${key}_stack"
func (l *lEntry) AddError(key string, val error) Entry {
	msg := val.Error()
	st := errors.ErrorStack(val)
	l.entry = l.entry.WithField(key, msg)
	l.entry = l.entry.WithField(key+"_stack", st)
	return l
}

// AddBool adds a bool value to the log statement.
func (l *lEntry) AddBool(key string, val bool) Entry {
	l.entry = l.entry.WithField(key, val)
	return l
}

// AddInt adds an integer value to the log statement.
func (l *lEntry) AddInt(key string, val int) Entry {
	l.entry = l.entry.WithField(key, val)
	return l
}

// AddStr adds a string value to the log statement.
func (l *lEntry) AddStr(key string, val string) Entry {
	l.entry = l.entry.WithField(key, val)
	return l
}

// AddTime adds a time value to the log statement.
func (l *lEntry) AddTime(key string, val time.Time) Entry {
	l.entry = l.entry.WithField(key, val)
	return l
}

// AddDur adds a duration value to the log statement.
func (l *lEntry) AddDur(key string, val time.Duration) Entry {
	l.entry = l.entry.WithField(key, val)
	return l
}

// AddAny adds any value to the log statement.
func (l *lEntry) AddAny(key string, val interface{}) Entry {
	l.entry = l.entry.WithField(key, val)
	return l
}

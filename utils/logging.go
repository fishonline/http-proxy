package utils

import (
	"fmt"
	"io"
	"math"
	"time"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	QUIET
)

var (
	logTimestampFormat = "Jan 02 15:04:05.000"
	processStart       = time.Now()

	NullLogger Logger = &NOPLogger{}
)

type LogLevel int

// Logger defines a simple logging interface
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	IsLevel(level LogLevel) bool
}

type NOPLogger struct{}

func (*NOPLogger) Debugf(format string, args ...interface{}) {}

func (*NOPLogger) Infof(format string, args ...interface{}) {}

func (*NOPLogger) Warningf(format string, args ...interface{}) {}

func (*NOPLogger) Errorf(format string, args ...interface{}) {}

func (l *NOPLogger) IsLevel(level LogLevel) bool {
	return false
}

// TimeLogger logs to stdout with a timestamp
type TimeLogger struct {
	debugW *io.Writer
	infoW  *io.Writer
	warnW  *io.Writer
	errorW *io.Writer
}

func NewTimeLogger(w *io.Writer, lvl LogLevel) *TimeLogger {
	l := &TimeLogger{}
	if lvl <= DEBUG {
		l.debugW = w
	}
	if lvl <= INFO {
		l.infoW = w
	}
	if lvl <= WARN {
		l.warnW = w
	}
	if lvl <= ERROR {
		l.errorW = w
	}
	return l
}

func (t *TimeLogger) Debugf(format string, args ...interface{}) {
	if t.debugW != nil {
		t.write(t.debugW, "DEBUG ", fmt.Sprintf(format, args...))
	}
}

func (t *TimeLogger) Infof(format string, args ...interface{}) {
	if t.infoW != nil {
		t.write(t.infoW, "INFO  ", fmt.Sprintf(format, args...))
	}
}

func (t *TimeLogger) Warningf(format string, args ...interface{}) {
	if t.warnW != nil {
		t.write(t.warnW, "WARN  ", fmt.Sprintf(format, args...))
	}
}

func (t *TimeLogger) Errorf(format string, args ...interface{}) {
	if t.errorW != nil {
		t.write(t.errorW, "ERROR ", fmt.Sprintf(format, args...))
	}
}

func (t *TimeLogger) write(w *io.Writer, prefix, str string) (n int, err error) {
	ts := time.Now()
	runningSecs := ts.Sub(processStart).Seconds()
	secs := int(math.Mod(runningSecs, 60))
	mins := int(runningSecs / 60)
	return fmt.Fprintf(*w, "%s%s - %dm%ds %s", prefix,
		ts.In(time.UTC).Format(logTimestampFormat),
		mins, secs, str)
}

func (t *TimeLogger) IsLevel(level LogLevel) bool {
	switch level {
	case DEBUG:
		return t.debugW != nil
	case INFO:
		return t.infoW != nil
	case WARN:
		return t.warnW != nil
	case ERROR:
		return t.errorW != nil
	}
	return false
}

package log

import (
	"errors"
	"fmt"
	"time"

	logrus "github.com/Sirupsen/logrus"
	stack "github.com/go-stack/stack"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	lfshook "github.com/rifflock/lfshook"
)

// START FROM go-kit/log

// A Valuer generates a log value. When passed to With or WithPrefix in a
// value element (odd indexes), it represents a dynamic value which is re-
// evaluated with each log event.
type Valuer func() interface{}

// bindValues replaces all value elements (odd indexes) containing a Valuer
// with their generated value.
func bindValues(keyvals []interface{}) {
	for i := 1; i < len(keyvals); i += 2 {
		if v, ok := keyvals[i].(Valuer); ok {
			keyvals[i] = v()
		}
	}
}

// Caller returns a Valuer that returns a file and line from a specified depth
// in the callstack. Users will probably want to use DefaultCaller.
func Caller(depth int) Valuer {
	return func() interface{} { return stack.Caller(depth) }
}

// DefaultCaller is a Valuer that returns the file and line where the Log
// method was invoked. It can only be used with log.With.
var DefaultCaller = Caller(3)

// ErrMissingValue is appended to keyvals slices with odd length to substitute
// the missing value.
var ErrMissingValue = errors.New("(MISSING)")

// END FROM go-kit/log

type eglooLogger struct {
	errorLevel logrus.Level
	logger     *logrus.Entry
}

type EglooLogger interface {
	Error(err error, msg string)
	Log(args ...interface{})
	Info(args ...interface{})
	Kinfo(args ...interface{})
	Warn(args ...interface{})
	Fatal(args ...interface{})
	WithFields(fields Fields) *Entry
	LogFatalIfError(err error, msg string)
}

type Entry struct {
	logrus.Entry
}

type Fields struct {
	logrus.Fields
}

func (l *eglooLogger) Error(err error, msg string) {
	l.logger.Error(err.Error())
}

func (l *eglooLogger) Log(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *eglooLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *eglooLogger) Kinfo(args ...interface{}) {
	stdLevel := logrus.StandardLogger().Level
	loggerLevel := l.logger.Logger.Level
	SetLevel(logrus.InfoLevel)
	l.logger.Logger.Level = logrus.InfoLevel

	l.logger.Info(args...)

	SetLevel(stdLevel)
	l.logger.Logger.Level = loggerLevel
}

func (l *eglooLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *eglooLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *eglooLogger) WithFields(fields Fields) *Entry {
	return &Entry{
		*l.logger.WithFields(fields.Fields)}
}

func (l *eglooLogger) LogFatalIfError(err error, msg string) {
	if err != nil {
		logrus.Fatalf("%s: %s", msg, err)
		l.logger.Fatalf("%s: %s", msg, err)
	}
}

func SetLevel(level logrus.Level) {
	logrus.SetLevel(level)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

func WithFields(fields Fields) *Entry {
	return &Entry{
		*logrus.WithFields(fields.Fields)}
}

func NewFields(args ...interface{}) Fields {
	fields := Fields{logrus.Fields{}}
	if len(args) == 0 {
		return fields
	}
	if len(args)%2 != 0 {
		args = append(args, ErrMissingValue)
	}

	for i := 0; i < len(args); i += 2 {
		fields.Fields[args[i].(string)] = args[i+1]
	}

	return fields
}

func NewEglooLogger(serviceName string) EglooLogger {
	writer, err := rotatelogs.New(
		fmt.Sprintf("/var/log/%s.log.%s", serviceName, "%Y%m%d%H%M"),
		rotatelogs.WithLinkName(fmt.Sprintf("/var/log/%s.log", serviceName)),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"serviceName": serviceName,
			"err":         err,
		}).Fatal("Failed to create rotatelogs")
		return nil
	}

	logger := logrus.WithFields(Fields{}.Fields)
	logger.Logger.Formatter = &logrus.JSONFormatter{}
	logger.Logger.Hooks.Add(lfshook.NewHook(lfshook.WriteMap{
		logrus.InfoLevel:  writer,
		logrus.ErrorLevel: writer}))

	return &eglooLogger{
		logger: logger}
}

func NewEglooLoggerWithFields(serviceName string, fields Fields) EglooLogger {
	writer, err := rotatelogs.New(
		fmt.Sprintf("/var/log/%s.log.%s", serviceName, "%Y%m%d%H%M"),
		rotatelogs.WithLinkName(fmt.Sprintf("/var/log/%s.log", serviceName)),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"serviceName": serviceName,
			"err":         err,
		}).Fatal("Failed to create rotatelogs")
		return nil
	}

	logger := logrus.WithFields(fields)
	logger.Logger.Formatter = &logrus.JSONFormatter{}
	logger.Logger.Hooks.Add(lfshook.NewHook(lfshook.WriteMap{
		logrus.InfoLevel:  writer,
		logrus.ErrorLevel: writer}))

	return &eglooLogger{
		logger: logger}
}

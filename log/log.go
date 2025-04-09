// Package log direct use logger
package log

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/natefinch/lumberjack"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	otellog "go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger *Logger
	loggers       = make(map[string]*Logger, 0)
	mu            sync.Mutex
)

var (
	defaultLogOpt = LoggerOpt{
		LogLevel:         zapcore.InfoLevel,
		Directory:        ".",
		TraceIDEnable:    true,
		MaxSize:          15,
		MaxBackups:       5,
		MaxAge:           365,
		IsDefault:        true,
		ConsoleLogEnable: true,
		EnableCaller:     true,
		LoggerProvider:   nil,
	}

	// CommonLogOpt use to easliy construct custom log options
	CommonLogOpt = defaultLogOpt

	logJsonEncodeCfg = zapcore.EncoderConfig{
		MessageKey:    "msg",                          // default msg key
		LevelKey:      "level",                        // log level key
		CallerKey:     "caller",                       // caller key
		TimeKey:       "time",                         // log time key
		StacktraceKey: "stack",                        // stack trace key
		LineEnding:    zapcore.DefaultLineEnding,      // log ends with "\n"
		EncodeLevel:   zapcore.LowercaseLevelEncoder,  // log level format "info"
		EncodeTime:    zapcore.RFC3339NanoTimeEncoder, // log time format
		EncodeCaller:  zapcore.FullCallerEncoder,      // Full caller path
	}

	logConsoleEncodeCfg = zapcore.EncoderConfig{
		MessageKey:    "msg",                            // default msg key
		LevelKey:      "level",                          // log level key
		CallerKey:     "caller",                         // caller key
		TimeKey:       "time",                           // log time key
		StacktraceKey: "stack",                          // stack trace key
		LineEnding:    zapcore.DefaultLineEnding,        // log ends with "\n"
		EncodeLevel:   zapcore.CapitalColorLevelEncoder, // log Level with colors
		EncodeTime:    zapcore.RFC3339TimeEncoder,       // not precise time
		EncodeCaller:  zapcore.ShortCallerEncoder,       // short caller path
	}
)

func newZapLogger(opt *LoggerOpt) *zap.Logger {
	// create log file folder
	if err := opt.CreateDirectory(); err != nil {
		panic(err)
	}

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   opt.GetLogFilePath(),
		MaxSize:    opt.MaxSize,
		MaxBackups: opt.MaxBackups,
		MaxAge:     opt.MaxAge,
	})

	cores := make([]zapcore.Core, 0)
	opts := make([]zap.Option, 0)

	jsonCore := zapcore.NewCore(zapcore.NewJSONEncoder(logJsonEncodeCfg), w, zap.DebugLevel)
	consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(logConsoleEncodeCfg), zapcore.Lock(os.Stdout), zap.DebugLevel)

	cores = append(cores, jsonCore)
	if opt.ConsoleLogEnable {
		cores = append(cores, consoleCore)
	}

	if opt.LoggerProvider != nil { // support OTLP
		otelzapCore := otelzap.NewCore(
			opt.Name,
			otelzap.WithLoggerProvider(opt.LoggerProvider),
		)
		cores = append(cores, otelzapCore)
	}

	opts = append(opts, zap.AddStacktrace(zap.ErrorLevel))
	if opt.EnableCaller {
		opts = append(opts, zap.AddCallerSkip(2))
		opts = append(opts, zap.AddCaller())
	}

	return zap.New(
		zapcore.NewTee(
			cores...,
		),
		opts...,
	)
}

// newLogger return a well-configed logger
func newLogger(name string, opt *LoggerOpt) *Logger {
	opt.Name = fmt.Sprintf("%s.log", name)
	logger := &Logger{
		opt: opt,
	}
	logger.zaplog = newZapLogger(logger.opt)
	logger.zaplog = logger.zaplog.With(zap.String("name", name))
	if opt.IsDefault {
		defaultLogger = logger
	}
	return logger
}

// GetLogger get logger with specified name
// or new a logger with options
func GetLogger(name string, opt *LoggerOpt) *Logger {
	if opt == nil {
		opt = &defaultLogOpt
	}
	mu.Lock()
	defer mu.Unlock()
	l, ok := loggers[name]
	if ok {
		return l
	}
	l = newLogger(name, opt)
	loggers[name] = l
	return l
}

// GetDefaultLogger get default logger
// if no default logger was initilized ** nil ** will be return
func GetDefaultLogger() *Logger {
	return defaultLogger
}

// Logger self defined Logger
type Logger struct {
	zaplog *zap.Logger
	opt    *LoggerOpt
}

// LoggerOpt configures the logger
type LoggerOpt struct {
	LogLevel         zapcore.Level
	Directory        string                  // log file directory
	Name             string                  // log file name
	TraceIDEnable    bool                    // enable traceid field
	MaxSize          int                     // Log File Max Size MB
	MaxBackups       int                     // The number of backup log file
	MaxAge           int                     // The days the log will be kept
	IsDefault        bool                    // is defalut logger?
	ConsoleLogEnable bool                    // enable console log?
	EnableCaller     bool                    // enable Caller?
	LoggerProvider   *otellog.LoggerProvider // when not nil, use otelzap bridge
}

// GetLogFilePath get log dst file path
func (opt LoggerOpt) GetLogFilePath() string {
	absPath, err := filepath.Abs(opt.Directory)
	if err != nil {
		panic(err)
	}
	return filepath.Join(absPath, opt.Name)
}

// CreateDirectory create logfile directory
func (opt LoggerOpt) CreateDirectory() error {
	if filepath.Clean(opt.Directory) == "" {
		return fmt.Errorf("directory: %s invalid", opt.Directory)
	}
	if _, err := os.Stat(opt.Directory); err == nil { // directory already exists
		return nil
	}
	err := os.MkdirAll(opt.Directory, os.ModePerm)
	return err
}

// WithDirectory setting log file directory
func (opt LoggerOpt) WithDirectory(path string) LoggerOpt {
	opt.Directory = path
	return opt
}

// WithTraceIDEnable setting traceid
func (opt LoggerOpt) WithTraceIDEnable(enable bool) LoggerOpt {
	opt.TraceIDEnable = enable
	return opt
}

// WithLogRetention sets log retention
func (opt LoggerOpt) WithLogRetention(maxSize int, maxBackups int, maxAge int) LoggerOpt {
	opt.MaxAge = maxAge
	opt.MaxBackups = maxBackups
	opt.MaxSize = maxSize
	return opt
}

// WithLogLevel sets log level
func (opt LoggerOpt) WithLogLevel(level zapcore.Level) LoggerOpt {
	opt.LogLevel = level
	return opt
}

// WithConsoleLog enable console log
func (opt LoggerOpt) WithConsoleLog(enable bool) LoggerOpt {
	opt.ConsoleLogEnable = enable
	return opt
}

// WithCaller enable caller info
func (opt LoggerOpt) WithCaller(enable bool) LoggerOpt {
	opt.EnableCaller = enable
	return opt
}

func (opt LoggerOpt) WithOtelLoggerProvider(provider *otellog.LoggerProvider) LoggerOpt {
	opt.LoggerProvider = provider
	return opt
}

// Debug the msg
func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.log(ctx, zapcore.DebugLevel, l.zaplog.Debug, msg, fields...)
}

// Info the msg
func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.log(ctx, zapcore.InfoLevel, l.zaplog.Info, msg, fields...)
}

// Warn the msg
func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.log(ctx, zapcore.WarnLevel, l.zaplog.Warn, msg, fields...)
}

// Error the msg
func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.log(ctx, zapcore.ErrorLevel, l.zaplog.Error, msg, fields...)
}

// Fatal the msg and exit with errcode 1
func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.log(ctx, zapcore.FatalLevel, l.zaplog.Fatal, msg, fields...)
}

func (l *Logger) log(
	ctx context.Context,
	logLevel zapcore.Level,
	logFunc func(msg string, fields ...zap.Field),
	msg string,
	fields ...zap.Field,
) {
	if logLevel < l.opt.LogLevel {
		return
	}
	var dst []zapcore.Field
	if l.opt.TraceIDEnable {
		dst = append(dst, zap.String("traceid", GetTraceIdWithCtx(ctx)))
	}
	// add remaining fields
	dst = append(dst, fields...)
	logFunc(msg, dst...)
}

// SetLevel setting log level
func (l *Logger) SetLevel(level zapcore.Level) {
	l.opt.LogLevel = level
}

// WithLoggerMetaFields set meta fields for logger
func (l *Logger) WithLoggerMetaFields(fields ...zapcore.Field) *Logger {
	l.zaplog = l.zaplog.With(fields...)
	return l
}

// NewFromLogger new a logger from a logger
func NewFromLogger(logger *Logger) *Logger {
	return &Logger{
		zaplog: logger.zaplog,
		opt:    logger.opt,
	}
}

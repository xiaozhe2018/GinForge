package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
	Fatal(msg string, fields ...any)
	With(fields ...zap.Field) Logger
	Desugar() *zap.Logger
}

type zapLogger struct{ l *zap.Logger }

func (z *zapLogger) Debug(msg string, fields ...any) { z.l.Sugar().Debugw(msg, fields...) }
func (z *zapLogger) Info(msg string, fields ...any)  { z.l.Sugar().Infow(msg, fields...) }
func (z *zapLogger) Warn(msg string, fields ...any)  { z.l.Sugar().Warnw(msg, fields...) }
func (z *zapLogger) Error(msg string, fields ...any) { z.l.Sugar().Errorw(msg, fields...) }
func (z *zapLogger) Fatal(msg string, fields ...any) { z.l.Sugar().Fatalw(msg, fields...) }
func (z *zapLogger) With(fields ...zap.Field) Logger { return &zapLogger{l: z.l.With(fields...)} }
func (z *zapLogger) Desugar() *zap.Logger            { return z.l }

func New(serviceName, logLevel string) Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder

	var level zapcore.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.DebugLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		level,
	)
	l := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.Fields(zap.String("service", serviceName)))
	return &zapLogger{l: l}
}

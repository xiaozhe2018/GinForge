package logger

import (
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
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

// New 创建日志器
// serviceName: 服务名称（用于生成日志文件名）
// logLevel: 日志级别 (debug/info/warn/error)
// logOutput: 日志输出方式 (stdout/file/both)
// logDir: 日志目录（当 logOutput 为 file 或 both 时使用）
func New(serviceName, logLevel, logOutput, logDir string) Logger {
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

	// 创建多个输出目标
	var cores []zapcore.Core

	// 标准输出（控制台）
	if logOutput == "stdout" || logOutput == "both" {
		stdoutCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, stdoutCore)
	}

	// 文件输出
	if logOutput == "file" || logOutput == "both" {
		// 确保日志目录存在
		if logDir == "" {
			logDir = "logs"
		}
		if err := os.MkdirAll(logDir, 0755); err != nil {
			// 如果创建目录失败，降级到 stdout
			stdoutCore := zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderCfg),
				zapcore.AddSync(os.Stdout),
				level,
			)
			cores = append(cores, stdoutCore)
		} else {
			// 每个服务一个日志文件：logs/{service-name}.log
			logFile := filepath.Join(logDir, serviceName+".log")
			fileWriter := &lumberjack.Logger{
				Filename:   logFile,
				MaxSize:    100, // MB
				MaxBackups: 10,
				MaxAge:     30, // days
				Compress:   true,
			}

			fileCore := zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderCfg),
				zapcore.AddSync(fileWriter),
				level,
			)
			cores = append(cores, fileCore)
		}
	}

	// 如果没有配置任何输出，默认使用 stdout
	if len(cores) == 0 {
		stdoutCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, stdoutCore)
	}

	// 合并多个 core
	core := zapcore.NewTee(cores...)
	l := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.Fields(zap.String("service", serviceName)))
	return &zapLogger{l: l}
}

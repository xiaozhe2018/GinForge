package event

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"
	
	"goweb/pkg/logger"
)

// LoggingMiddleware 日志中间件
func LoggingMiddleware(log logger.Logger) Middleware {
	return func(next Handler) Handler {
		return func(e Event) error {
			startTime := time.Now()
			
			err := next(e)
			
			duration := time.Since(startTime)
			
			if err != nil {
				log.Error("event handler error",
					err,
					"event_type", e.Type(),
					"duration", duration.String())
			} else {
				log.Debug("event handled",
					"event_type", e.Type(),
					"duration", duration.String())
			}
			
			return err
		}
	}
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware(log logger.Logger) Middleware {
	return func(next Handler) Handler {
		return func(e Event) (err error) {
			defer func() {
				if r := recover(); r != nil {
					stack := string(debug.Stack())
					log.Error("event handler panic recovered",
						fmt.Errorf("%v", r),
						"event_type", e.Type(),
						"stack", stack)
					err = fmt.Errorf("panic: %v", r)
				}
			}()
			
			return next(e)
		}
	}
}

// TimeoutMiddleware 超时中间件
func TimeoutMiddleware(timeout time.Duration) Middleware {
	return func(next Handler) Handler {
		return func(e Event) error {
			ctx := e.Context()
			if ctx == nil {
				ctx = context.Background()
			}
			
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			
			// 创建带有新上下文的事件
			eventWithTimeout := NewEventWithContext(ctx, e.Type(), e.Data())
			
			done := make(chan error, 1)
			go func() {
				done <- next(eventWithTimeout)
			}()
			
			select {
			case err := <-done:
				return err
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

// RetryMiddleware 重试中间件
func RetryMiddleware(maxRetries int, delay time.Duration) Middleware {
	return func(next Handler) Handler {
		return func(e Event) error {
			var err error
			
			for i := 0; i <= maxRetries; i++ {
				err = next(e)
				if err == nil {
					return nil
				}
				
				if i < maxRetries {
					time.Sleep(delay)
				}
			}
			
			return err
		}
	}
}

// ValidationMiddleware 验证中间件
func ValidationMiddleware(validator func(Event) error) Middleware {
	return func(next Handler) Handler {
		return func(e Event) error {
			if err := validator(e); err != nil {
				return err
			}
			
			return next(e)
		}
	}
}

// MetricsMiddleware 指标中间件
func MetricsMiddleware(recorder func(eventType string, duration time.Duration, err error)) Middleware {
	return func(next Handler) Handler {
		return func(e Event) error {
			startTime := time.Now()
			
			err := next(e)
			
			duration := time.Since(startTime)
			recorder(e.Type(), duration, err)
			
			return err
		}
	}
}

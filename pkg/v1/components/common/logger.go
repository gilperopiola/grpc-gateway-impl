package common

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

/* ----------------------------------- */
/*             - Logger -              */
/* ----------------------------------- */

type Logger struct {
	*zap.Logger
	sugar *zap.SugaredLogger
}

// NewLogger returns a new instance of Logger.
func NewLogger(isProd bool, opts ...zap.Option) *Logger {
	newLoggerFunc := zap.NewDevelopment
	if isProd {
		newLoggerFunc = zap.NewProduction
	}

	logger, err := newLoggerFunc(opts...)
	if err != nil {
		log.Fatalf(errs.FatalErrMsgCreatingLogger, err)
	}

	return &Logger{Logger: logger, sugar: logger.Sugar()}
}

// NewLoggerOptions returns the default options for the Logger.
func NewLoggerOptions() []zap.Option {
	return []zap.Option{
		zap.AddStacktrace(zap.DPanicLevel),
	}
}

// LogGRPC logs the gRPC Request's info when it finishes executing.
func (l *Logger) LogGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)

	if err != nil {
		l.sugar.Errorw("gRPC Error", endpointField(info.FullMethod), durationField(duration), errorField(err))
	} else {
		l.sugar.Infow("gRPC Request", endpointField(info.FullMethod), durationField(duration))
	}

	return resp, err
}

// LogHTTP logs the HTTP Request's info when it finishes executing.
func (l *Logger) LogHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		l.sugar.Infow("HTTP Request", endpointField(r.Method+" "+r.URL.Path), durationField(duration))

		// Most HTTP logs come with a gRPC log before, as HTTP acts as a gateway to gRPC.
		// As such, we add a new line to separate the logs and easily identify different requests.
		// The only exception would be if there was an error before calling the gRPC handlers.
		l.sugar.Infoln("")
	})
}

func endpointField(value string) zap.Field {
	return zap.String("endpoint", value)
}

func durationField(value time.Duration) zap.Field {
	return zap.Duration("duration", value)
}

func errorField(err error) zap.Field {
	return zap.Error(err)
}

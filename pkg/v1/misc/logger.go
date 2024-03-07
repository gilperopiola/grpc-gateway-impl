package misc

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

// NewLogger returns a new instance of *zap.Logger.
func NewLogger(isProd bool, opts []zap.Option) *zap.Logger {
	newLoggerFunc := zap.NewDevelopment
	if isProd {
		newLoggerFunc = zap.NewProduction
	}

	logger, err := newLoggerFunc(opts...)
	if err != nil {
		log.Fatalf(errs.FatalErrMsgCreatingLogger, err)
	}

	return logger
}

// NewLoggerOptions returns the default options for the logger.
// For now it only adds a stack trace to panic logs.
func NewLoggerOptions() []zap.Option {
	return []zap.Option{
		zap.AddStacktrace(zap.DPanicLevel),
	}
}

// LogGRPC logs the gRPC request's info when it finishes executing.
func LogGRPC(logger *zap.Logger) grpc.UnaryServerInterceptor {
	sugar := logger.Sugar()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			sugar.Errorw("gRPC Error", endpointField(info.FullMethod), durationField(duration), errorField(err))
		} else {
			sugar.Infow("gRPC Request", endpointField(info.FullMethod), durationField(duration))
		}

		// After logging the request, we return the response and error because this is an interceptor.
		return resp, err
	}
}

// LogHTTP logs the HTTP Request's info when it finishes executing.
func LogHTTP(next http.Handler, logger *zap.Logger) http.Handler {
	sugar := logger.Sugar()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		sugar.Infow("HTTP Request", endpointField(r.Method+" "+r.URL.Path), durationField(duration))

		// Most HTTP logs come with a gRPC log before, as HTTP acts as a gateway to gRPC.
		// As such, we add a new line to separate the logs and easily identify different requests.
		// The only exception would be if there was an error before calling the gRPC handlers.
		sugar.Infoln("")
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

package log

import (
	"fmt"
	"os"

	stdlog "log"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
)

type Logger interface {
	Log(keyvals ...interface{}) error
}

func NewLogger() Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	stdlog.SetOutput(log.NewStdlibAdapter(logger))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	// We're wrapping gotkit, so DefaultCaller would show: log.go:23
	logger = log.With(logger, "caller", log.Caller(3))

	return logger
}

// NewNopLogger is go-kit/log.NewNopLogger
func NewNopLogger() log.Logger {
	return log.NewNopLogger()
}

// NewSyncLogger is useful for debugging, use sparingly
func NewSyncLogger() log.Logger {
	return log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
}

// Str is used to log unknown types using fmt
func Str(i any) string { return fmt.Sprintf("%+v", i) }

// GinFormatterWithUTCAndBodySize is the default gin loggger with:
//   - UTC times
//   - logs the response size in bytes
func GinFormatterWithUTCAndBodySize(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}
	return fmt.Sprintf("[DUNE] %v |%s %3d %s| %13v | %6v bytes | %15s |%s %-7s %s %#v\n%s",
		param.TimeStamp.UTC().Format("2006/01/02 15:04:05.000"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.BodySize,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

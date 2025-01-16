package observability

import (
	"fmt"
	"os"

	"go.uber.org/zap/zapcore"
)

func ZapAdapter() func(core zapcore.Core) zapcore.Core {
	return func(core zapcore.Core) zapcore.Core {
		return zapcore.RegisterHooks(core, func(entry zapcore.Entry) error {
			if entry.Level == zapcore.ErrorLevel {
				hostname, _ := os.Hostname()

				fmt.Printf("ERROR\nHostname: %s\nENV: %s\nExtra: %s\nMessage: %s\nTimestamp: %s\nLogger: %s", hostname, "production", map[string]any{
					"caller": entry.Caller.String(),
					"stack":  entry.Stack,
				}, entry.Message, entry.Time, entry.LoggerName)
			}

			return nil
		})
	}
}

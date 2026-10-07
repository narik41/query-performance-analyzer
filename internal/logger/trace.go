package logger

import "database-query-analyzer/internal/config"

func GetTraceKey() string {
	return defaultString(config.Default().GetString("logger.traceKey"), "debugId")
}

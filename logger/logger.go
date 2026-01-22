package logger

//go:generate mockgen -source=logger.go -destination=mock_logger.go -package=logger
type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
	Debug(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Close()
}

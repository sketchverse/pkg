package logger

import "testing"

func init() {
	InitLogger(nil)
}

func TestLogger(t *testing.T) {
	logger := Logger()
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")
	logger.Debug("This is a debug message")
	// logger.Fatal("This is a fatal message")

}

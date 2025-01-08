package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Logger wraps the standard log package functionality
type Logger struct {
	file *os.File
	logger *log.Logger
	mu sync.Mutex
	filepath string
}

// Creates a new Logger instance that writes to both stdout and a file
func New(filename string) (*Logger, error) {
	if filename == "" {
		filename = "application.log"
	}

	// Get the current working director
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Failed to get working directory: %v", err)
	}

	logsDir := filepath.Join(wd, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("Error creating logs directory: %v", err)
	}

	// Create date-based directory
	currentTime := time.Now()
	dateDir := currentTime.Format("2006-01-02") // YYYY-MM-DD
	datePath := filepath.Join(logsDir, dateDir)
	if err := os.MkdirAll(datePath, 0755); err != nil {
		return nil, fmt.Errorf("Failed to create date directory: %v", err)
	}

	// create full filepath
	fullPath := filepath.Join(datePath, filename)

	// open log file with append mode, create if doesnt exist
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Failed to open log file: %v", err)
	}

	// create multi-writer to write to both file and stdout
	logger := log.New(file, "", log.LstdFlags)

	return &Logger{
		file: file,
		logger: logger,
		filepath: fullPath,
	}, nil
}

// closes the log file
func (l *Logger) Close() error {
	return l.file.Close()
}

// returns the file and line number of the caller
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown:0"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// formats the log message with timestamp and caller info
func formatLog(format string, v ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	callerInfo := getCallerInfo()
	message := fmt.Sprintf(format, v...)
	return fmt.Sprintf("[%s] [%s] %s", timestamp, callerInfo, message)
}

// logs a message
func (l *Logger) Print(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	msg := formatLog("%v", v...)
	l.logger.Print(msg)
	fmt.Println(msg)
}

// logs a formatted message
func (l *Logger) Printf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	msg := formatLog(format, v...)
	l.logger.Print(msg)
	fmt.Println(msg)
}

// logs a message and exits with status 1
func (l *Logger) Fatal(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	msg := formatLog("%v", v...)
	l.logger.Print(msg)
	fmt.Println(msg)
	os.Exit(1)
}

// logs a formatted message and exits with a status 1
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	msg := formatLog(format, v...)
	l.logger.Print(msg)
	fmt.Println(msg)
	os.Exit(1)
}

// formats using the default formats for its operands and returns the string
func (l *Logger) Sprint(v ...interface{}) string {
	return fmt.Sprint(v...)
}

// formats according to a format specifier and returns the string
func (l *Logger) Sprintf(format string, v ...interface{}) string {
	return fmt.Sprintf(format, v...)
}

// logs an error message
func (l *Logger) Error(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	msg := formatLog("[ERROR] %v", v...)
	l.logger.Print(msg)
	fmt.Println(msg)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	userMsg := fmt.Sprintf(format, v...)
	msg := formatLog("[ERROR] %s", userMsg)
	l.logger.Print(msg)
	fmt.Println(msg)
}
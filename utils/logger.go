package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	// File loggers
	fileInfoLog     *log.Logger
	fileErrorLog    *log.Logger
	fileFatalLog    *log.Logger
	fileRequestLog  *log.Logger
	fileResponseLog *log.Logger

	// Console loggers
	consoleInfoLog     *log.Logger
	consoleErrorLog    *log.Logger
	consoleFatalLog    *log.Logger
	consoleRequestLog  *log.Logger
	consoleResponseLog *log.Logger

	logFile *os.File

	// ANSI Color Codes
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
)

func InitLogger() {
	// Konfigurasi folder log dan nama file log berdasarkan tanggal
	logFolder := "logs"
	currentDate := time.Now().Format("2006-01-02") // Format tanggal: YYYY-MM-DD
	logFileName := fmt.Sprintf("%s/%s.log", logFolder, currentDate)

	// Membuat folder jika belum ada
	if _, err := os.Stat(logFolder); os.IsNotExist(err) {
		err := os.MkdirAll(logFolder, 0755)
		if err != nil {
			log.Fatalf("Failed to create log directory: %v", err)
		}
	}

	// Membuka file log
	var err error
	logFile, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Set up File loggers (Plain)
	fileInfoLog = log.New(logFile, "INFO: ", log.Ldate|log.Ltime)
	fileErrorLog = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime)
	fileFatalLog = log.New(logFile, "FATAL: ", log.Ldate|log.Ltime)
	fileRequestLog = log.New(logFile, "REQUEST: ", log.Ldate|log.Ltime)
	fileResponseLog = log.New(logFile, "RESPONSE: ", log.Ldate|log.Ltime)

	// Set up Console loggers (Colorful)
	consoleInfoLog = log.New(os.Stdout, Blue+"INFO: "+Reset, log.Ldate|log.Ltime)
	consoleErrorLog = log.New(os.Stdout, Red+"ERROR: "+Reset, log.Ldate|log.Ltime)
	consoleFatalLog = log.New(os.Stdout, Red+"FATAL: "+Reset, log.Ldate|log.Ltime)
	consoleRequestLog = log.New(os.Stdout, Cyan+"REQUEST: "+Reset, log.Ldate|log.Ltime)
	consoleResponseLog = log.New(os.Stdout, Magenta+"RESPONSE: "+Reset, log.Ldate|log.Ltime)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var bodyBytes []byte
		var parsedHeader []byte

		parsedHeader, err := json.Marshal(r.Header)
		if err != nil {
			Error("Failed to read request header for logging: %v", err)
		}

		if r.Body != nil {
			bodyBytes, err = io.ReadAll(r.Body)
			if err != nil {
				Error("Failed to read request body for logging: %v", err)
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		requestId := GetRequestID(r.Context())

		Request("%s %s | BODY=%s | HEADER=%s | [RequestId=%s]", r.Method, r.URL.Path, string(bodyBytes), string(parsedHeader), requestId)

		next.ServeHTTP(w, r)
	})
}

func Info(msg string, args ...any) {
	fileInfoLog.Printf(msg, args...)
	consoleInfoLog.Printf(msg, args...)
}

func Error(msg string, args ...any) {
	fileErrorLog.Printf(msg, args...)
	consoleErrorLog.Printf(msg, args...)
}

func Fatal(msg string, args ...any) {
	fileFatalLog.Printf(msg, args...)
	consoleFatalLog.Printf(msg, args...)
	if logFile != nil {
		logFile.Close()
	}
	os.Exit(1)
}

func Request(msg string, args ...any) {
	fileRequestLog.Printf(msg, args...)
	consoleRequestLog.Printf(msg, args...)
}

func Response(msg string, args ...any) {
	fileResponseLog.Printf(msg, args...)
	consoleResponseLog.Printf(msg, args...)
}

// Ensure to call this function to properly close the log file before application exit.
func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

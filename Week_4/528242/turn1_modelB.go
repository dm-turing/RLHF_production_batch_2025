package main  
import (  
    "fmt"
    "math/rand"
    "os"
    "strings"
    "sync"
    "time"

    "github.com/sirupsen/logrus"
)

// SeverityLevel represents the severity of an error message.
type SeverityLevel int

const (
    // DEBUG severity level
    DEBUG SeverityLevel = iota
    // INFO severity level
    INFO
    // WARNING severity level
    WARNING
    // ERROR severity level
    ERROR
    // CRITICAL severity level
    CRITICAL
    // FATAL severity level
    FATAL
)

// severityLevelToString maps SeverityLevel to its corresponding string representation.
var severityLevelToString = map[SeverityLevel]string{
    DEBUG:    "DEBUG",
    INFO:     "INFO",
    WARNING:  "WARNING",
    ERROR:    "ERROR",
    CRITICAL: "CRITICAL",
    FATAL:    "FATAL",
}

// ErrorMessage represents an error message with a severity level.
type ErrorMessage struct {
    Message    string
    Severity   SeverityLevel
    Metadata   map[string]string
}

func (em ErrorMessage) String() string {
    var metadataStr string
    for k, v := range em.Metadata {
        metadataStr += fmt.Sprintf(" %s=%s", k, v)
    }
    return fmt.Sprintf("%s%s: %s", severityLevelToString[em.Severity], metadataStr, em.Message)
}

// GenerateRandomError generates a random error message with a severity level.
func GenerateRandomError() ErrorMessage {
    severity := SeverityLevel(rand.Intn(6))
    message := fmt.Sprintf("Random error message %d", rand.Intn(100))
    return ErrorMessage{
        Message:  message,
        Severity: severity,
        Metadata: map[string]string{
            "source":    "RandomSource",
            "timestamp": time.Now().Format(time.RFC3339),
        },
    }
}

// ComposeErrorMessage composes an error message using a callback function.
func ComposeErrorMessage(cb func(em *ErrorMessage)) ErrorMessage {
    em := GenerateRandomError()
    cb(&em)
    return em
}

// AddRandomMetadata is a callback function that adds random metadata to an error message.
func AddRandomMetadata(em *ErrorMessage) {
    em.Metadata["randomKey"] = fmt.Sprintf("randomValue%d", rand.Intn(100))
}

// Logger is a custom logger that handles logging based on severity levels.
type Logger struct {
    *logrus.Logger
    severityLevel SeverityLevel
    wg            sync.WaitGroup
    stop          chan struct{}
}

// NewLogger creates a new Logger instance.
func NewLogger(severityLevel SeverityLevel) *Logger {
    logger := &Logger{
        Logger:        logrus.New(),
        severityLevel: severityLevel,
        stop:          make(chan struct{}),
    }
    logger.Out = os.Stdout
    logger.Formatter = &logrus.TextFormatter{
        DisableTimestamp: true,
    }
    return logger
}

// Start starts the logging loop.
func (l *Logger) Start() {
    l.wg.Add(1)
    go func() {
        defer l.wg.Done()
        ticker := time.NewTicker(time.Second)
        defer ticker.Stop()
        for {
            select {
            case <-l.stop:
                return
            case <-ticker.C:
                em := ComposeErrorMessage(AddRandomMetadata)
                l.logMessage(em)
            }
package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// ErrorCode is an enum for error codes
type ErrorCode string

const (
	InvalidInput  ErrorCode = "InvalidInput"
	InternalError ErrorCode = "InternalError"
	UserNotFound  ErrorCode = "UserNotFound"
)

// ErrorContext holds error context details
type ErrorContext struct {
	Context   string
	UserID    int
	ErrorCode ErrorCode
}

// NewError creates a new error with context and error code
func NewError(ctx ErrorContext, msg string) error {
	return fmt.Errorf("context: %s, user id: %d, error code: %s, message: %s", ctx.Context, ctx.UserID, ctx.ErrorCode, msg)
}

// ChainError chains an error with the same context
func ChainError(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}

// LogError logs an error to the console
func LogError(err error) {
	log.Println(err.Error())
}

// GenerateRandomError generates a random error
func GenerateRandomError(ctx ErrorContext) error {
	errorCodes := []ErrorCode{InvalidInput, InternalError, UserNotFound}
	randomIndex := rand.Intn(len(errorCodes))
	ctx.ErrorCode = errorCodes[randomIndex]
	msg := fmt.Sprintf("Random error: %s", ctx.ErrorCode)
	return NewError(ctx, msg)
}

// Main function
func main() {
	rand.Seed(time.Now().UnixNano())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			c := ""
			for _, v := range rand.Perm(5) {
				c += strconv.Itoa(v)
			}
			ctx := ErrorContext{
				Context: c,
				UserID:  rand.Intn(1000),
			}
			err := GenerateRandomError(ctx)
			LogError(err)
			if rand.Intn(2) == 0 {
				err = ChainError(err, "Additional context message")
				LogError(err)
			}
			time.Sleep(time.Second)
		}
	}()

	wg.Wait()
}

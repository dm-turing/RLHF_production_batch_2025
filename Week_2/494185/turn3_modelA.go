package main

import (
	"context"
	"fmt"
	"time"

	"go-retry/retry"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type APIError struct {
	code int
	err  error
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error: status code %d, error %v", e.code, e.err)
}

func callAPI(ctx context.Context, url string) error {
	logger := zap.L.With(zap.String("url", url))
	restClient := resty.New().SetContext(ctx)

	res, err := restClient.R().Get(url)
	if err != nil {
		return err
	}

	if res.StatusCode() >= 400 {
		return &APIError{code: res.StatusCode(), err: xerrors.Errorf("bad response: %v", res.String())}
	}

	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	url := "https://invalid-url.com/api"
	logger := zap.New(zap.ErrorLevel())

	// Wrap with logger and use retry package
	retryOptions := retry.New(retry.Options{
		Attempts:    5,
		InitialWait: 100 * time.Millisecond,
		MaxWait:     1 * time.Second,
	})

	_, err := retry.Do(ctx, retryOptions, func() error {
		return callAPI(ctx, url)
	})

	if err != nil {
		logger.Error("callAPI failed", zap.Error(err))
	} else {
		logger.Info("callAPI succeeded")
	}
}

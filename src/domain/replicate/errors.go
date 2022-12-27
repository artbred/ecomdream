package replicate

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrorCudaOutOfMemory = errors.New("cuda out of memory")
)

func parseError(err error) error {
	if strings.Contains(err.Error(), "context canceled") {
		return context.Canceled
	}

	if strings.Contains(err.Error(), "context deadline exceeded") {
		return context.DeadlineExceeded
	}

	if strings.Contains(err.Error(), "CUDA out of memory") {
		return ErrorCudaOutOfMemory
	}

	return err
}


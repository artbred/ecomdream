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
	if errors.As(err, &context.Canceled) {
		return context.Canceled
	}

	if errors.As(err, &context.DeadlineExceeded) {
		return context.DeadlineExceeded
	}

	if strings.Contains(err.Error(), "CUDA out of memory") {
		return ErrorCudaOutOfMemory
	}

	return err
}


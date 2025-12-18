package launcher

import (
	"context"
	"errors"
)

var ErrUnsupported = errors.New("launcher not supported")

type Launcher interface {
	Open(ctx context.Context, target string) error
}

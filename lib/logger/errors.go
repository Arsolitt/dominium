package logger

import (
	"context"
	"errors"
)

type errorWithLogCtx struct {
	next error
	ctx  logData
}

func (e *errorWithLogCtx) Error() string {
	return e.next.Error()
}

func WrapError(ctx context.Context, err error) error {
	c := logData{}
	if x, ok := ctx.Value(dataKey).(logData); ok {
		c = x
	}
	return &errorWithLogCtx{
		next: err,
		ctx:  c,
	}
}

func ErrorCtx(ctx context.Context, err error) context.Context {
	var e *errorWithLogCtx
	if errors.As(err, &e) {
		return context.WithValue(ctx, dataKey, e.ctx)
	}
	return ctx
}

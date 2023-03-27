package decorator

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

type Handler[A any, R any] interface {
	Handle(ctx context.Context, args A) (R, error)
}

func WithLogging[A any, R any](handler Handler[A, R], logger *logrus.Entry) Handler[A, R] {
	return loggingHandler[A, R]{
		base:   handler,
		logger: logger,
	}
}

type loggingHandler[A any, R any] struct {
	base   Handler[A, R]
	logger *logrus.Entry
}

func (handler loggingHandler[A, R]) Handle(ctx context.Context, args A) (result R, err error) {
	logger := handler.logger.WithFields(
		logrus.Fields{
			"Handler": fmt.Sprintf("%T", handler.base),
			"Args":    fmt.Sprintf("%v", args),
		},
	)

	logger.Debug("Start")
	defer func() {
		if err == nil {
			logger.Info("Success")
		} else {
			logger.WithError(err).Error("Fail")
		}
	}()

	return handler.base.Handle(ctx, args)
}

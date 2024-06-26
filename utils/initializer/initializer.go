package initializer

import (
	"github.com/dmitriibb/go-common/logging"
	"sync"
)

type Initializer interface {
	Init(f func() error) error
	InitWithArgs(initFunc func(fArgs ...any) error, args ...interface{}) error
}

func New(logger logging.Logger) Initializer {
	init := initializer{
		initialized: false,
		logger:      logger,
	}
	return &init
}

type initializer struct {
	mu          sync.Mutex
	initialized bool
	logger      logging.Logger
}

func (in *initializer) Init(initFunc func() error) error {
	in.mu.Lock()
	defer in.mu.Unlock()
	if in.initialized {
		in.logger.Debug("already initialized")
		return nil
	}
	in.logger.Debug("initializing")
	err := initFunc()
	if err != nil {
		in.logger.Error("not initialized because '%v'", err.Error())
		return err
	}
	in.initialized = true
	in.logger.Debug("initialized")
	return nil
}

func (in *initializer) InitWithArgs(initFunc func(fArgs ...any) error, args ...interface{}) error {
	in.mu.Lock()
	defer in.mu.Unlock()
	if in.initialized {
		in.logger.Debug("already initialized")
		return nil
	}
	in.logger.Debug("initializing")
	err := initFunc(args...)
	if err != nil {
		in.logger.Error("not initialized because '%v'", err.Error())
		return err
	}
	in.initialized = true
	in.logger.Debug("initialized")
	return nil
}

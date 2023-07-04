package deamonizer

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Daemonizer struct {
	logger *zap.Logger
	//Never read or close this channel!
	errCh     chan error
	ctx       context.Context
	cancelCtx context.CancelFunc
	wg        sync.WaitGroup
}

const DEFAULT_TIMEOUT_SHUTDOWN = 5 * time.Second

func NewDaemonizer(logger *zap.Logger) Daemonizer {
	ctx, cancel := context.WithCancel(context.Background())
	return Daemonizer{
		logger:    logger,
		errCh:     make(chan error),
		ctx:       ctx,
		cancelCtx: cancel,
		wg:        sync.WaitGroup{},
	}
}

func (d *Daemonizer) Start() {
	canceling := false

	go func() {
		for err := range d.errCh {
			if !canceling {
				d.logger.Error("Interrupt by error", zap.Error(err))
				d.cancelCtx()
				canceling = true
			} else {
				d.logger.Error("Error while canceling", zap.Error(err))
			}
		}
	}()

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	select {
	case <-d.ctx.Done():
		d.logger.Debug("Context is over")
		if err := d.ctx.Err(); err != nil {
			d.logger.Error("Context interrupt by error", zap.Error(err))
		}
	case sig := <-done:
		d.logger.Info("Catch signal", zap.Any("signal", sig))
		canceling = true
		d.cancelCtx()
	}
}

func (d *Daemonizer) GracefulShutdown(timeout time.Duration) {

	c := make(chan struct{})

	go func() {
		defer close(c)
		d.wg.Wait()
	}()

	select {
	case <-c:
		d.logger.Error("Correct shutdown")
	case <-time.After(timeout):
		d.logger.Error("End by timeout")
	}
}

func (d *Daemonizer) AddDaemon(run func() error, shutdown func() error) {
	d.wg.Add(1)
	go func() {
		err := run()
		if err != nil {
			d.errCh <- err
		}
		d.wg.Done()
	}()

	select {
	case <-d.ctx.Done():
		if err := shutdown(); err != nil {
			d.logger.Error("Interrupt shutdown", zap.Error(err))
		}
		return
	}
}

package wsreply

import (
	"context"
	"sync"

	"github.com/sheirys/wsreply/broker"
)

type Application struct {
	Broker   broker.Broker
	ctx      context.Context
	stopFunc context.CancelFunc
	wg       *sync.WaitGroup
}

func (a *Application) Start() error {
	a.ctx, a.stopFunc = context.WithCancel(context.Background())
	a.wg = &sync.WaitGroup{}
	return nil
}

func (a *Application) Stop() error {
	a.stopFunc()
	a.wg.Wait()
	return nil
}

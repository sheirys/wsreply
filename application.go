package wsreply

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/sheirys/wsreply/broker"
)

type Application struct {
	ctx      context.Context
	stopFunc context.CancelFunc
	wg       *sync.WaitGroup
	http     *http.Server

	Addr   string
	Broker broker.Broker
	Log    *log.Logger
}

func (a *Application) Start() error {
	a.ctx, a.stopFunc = context.WithCancel(context.Background())
	a.wg = &sync.WaitGroup{}

	a.http = &http.Server{
		Addr:    a.Addr,
		Handler: a.router(),
	}
	a.Log.Printf("starting on %s", a.Addr)

	if err := a.Broker.Start(a.ctx, a.wg); err != nil {
		return err
	}

	a.wg.Add(1)
	go func() {
		if err := a.http.ListenAndServe(); err != nil {
			a.Log.Println(err)
		}
		a.wg.Done()
	}()

	return nil
}

func (a *Application) Stop() error {
	a.stopFunc()
	a.http.Shutdown(a.ctx)
	a.wg.Wait()
	return nil
}

package wsreply

import (
	"context"
	"net/http"
	"sync"

	"github.com/sheirys/wsreply/broker"
	"github.com/sirupsen/logrus"
)

type Application struct {
	ctx      context.Context
	stopFunc context.CancelFunc
	wg       *sync.WaitGroup
	http     *http.Server

	Addr   string
	Broker broker.Broker
	Log    *logrus.Logger
	Debug  bool
}

func (a *Application) Init() error {
	a.ctx, a.stopFunc = context.WithCancel(context.Background())
	a.wg = &sync.WaitGroup{}

	a.http = &http.Server{
		Addr:    a.Addr,
		Handler: a.router(),
	}
	return nil
}

func (a *Application) StartHTTP() error {
	a.Log.WithField("host", a.Addr).Info("starting server")
	go func() {
		a.wg.Add(1)
		if err := a.http.ListenAndServe(); err != nil {
			a.Log.Println(err)
		}
		a.wg.Done()
	}()
	return nil
}

func (a *Application) StartBroker() error {
	return a.Broker.Start()
}

func (a *Application) Stop() error {
	a.stopFunc()
	a.http.Shutdown(a.ctx)
	a.Broker.Stop()
	a.wg.Wait()
	return nil
}

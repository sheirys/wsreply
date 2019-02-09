package wsreply

import (
	"context"
	"net/http"
	"sync"

	"github.com/sheirys/wsreply/broker"
	"github.com/sirupsen/logrus"
)

type Server struct {
	ctx      context.Context
	stopFunc context.CancelFunc
	wg       *sync.WaitGroup
	http     *http.Server

	Addr   string
	Broker broker.Broker
	Log    *logrus.Logger
	Debug  bool
}

func (s *Server) Init() error {
	s.ctx, s.stopFunc = context.WithCancel(context.Background())
	s.wg = &sync.WaitGroup{}

	s.http = &http.Server{
		Addr:    s.Addr,
		Handler: s.router(),
	}
	return nil
}

func (s *Server) StartHTTP() error {
	s.Log.WithField("host", s.Addr).Info("starting server")
	go func() {
		s.wg.Add(1)
		if err := s.http.ListenAndServe(); err != nil {
			s.Log.Println(err)
		}
		s.wg.Done()
	}()
	return nil
}

func (s *Server) StartBroker() error {
	return s.Broker.Start()
}

func (s *Server) Stop() error {
	s.stopFunc()
	s.http.Shutdown(s.ctx)
	s.Broker.Stop()
	s.wg.Wait()
	return nil
}

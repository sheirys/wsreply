package wsreply

import (
	"context"
	"net/http"
	"sync"

	"github.com/sheirys/wsreply/broker"
	"github.com/sirupsen/logrus"
)

// Server is main wsreply application. Here broker and http server are hold
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

// Init should be called before using other Server functions. Here various
// variable initiations should be applied.
func (s *Server) Init() error {
	s.ctx, s.stopFunc = context.WithCancel(context.Background())
	s.wg = &sync.WaitGroup{}

	s.http = &http.Server{
		Addr:    s.Addr,
		Handler: s.router(),
	}
	return nil
}

// StartHTTP will start HTTP server on Server struct.
func (s *Server) StartHTTP() error {
	s.Log.WithField("host", s.Addr).Info("starting server")
	s.wg.Add(1)
	go func() {
		if err := s.http.ListenAndServe(); err != nil {
			s.Log.Println(err)
		}
		s.wg.Done()
	}()
	return nil
}

// StartBroker will start broker on Server struct.
func (s *Server) StartBroker() error {
	return s.Broker.Start()
}

// Stop will stop http and broker services.
func (s *Server) Stop() error {
	s.stopFunc()
	s.http.Shutdown(s.ctx)
	s.Broker.Stop()
	s.wg.Wait()
	return nil
}

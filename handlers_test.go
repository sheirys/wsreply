package wsreply_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/sheirys/wsreply"
	"github.com/sheirys/wsreply/broker"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testServer will initialize default wsreply server for testing.
func testServer() *wsreply.Server {
	srv := &wsreply.Server{
		Broker: broker.NewInMemBroker(true),
		Log:    logrus.New(),
	}
	srv.Init()
	srv.StartBroker()
	return srv
}

// connectWithWS will initialize ws connection to testserver
func connectWithWS(s *httptest.Server) (*websocket.Conn, error) {
	testURL := strings.Replace(s.URL, "http", "ws", 1)
	ws, _, err := websocket.DefaultDialer.Dial(testURL, nil)
	return ws, err
}

// TestWSSubscriberCanConnect will test if subscriber is able to connect to
// server.
func TestWSSubscriberCanConnect(t *testing.T) {
	srv := testServer()
	defer srv.Stop()

	s := httptest.NewServer(http.HandlerFunc(srv.WSSubscriber))
	defer s.Close()

	ws, err := connectWithWS(s)
	require.NoError(t, err)
	defer ws.Close()
}

func TestWSPublisherCanConnect(t *testing.T) {
	srv := testServer()
	defer srv.Stop()

	s := httptest.NewServer(http.HandlerFunc(srv.WSPublisher))
	defer s.Close()

	ws, err := connectWithWS(s)
	require.NoError(t, err)
	defer ws.Close()
}

func TestWSPublisherCanPublish(t *testing.T) {
	srv := testServer()
	defer srv.Stop()

	s := httptest.NewServer(http.HandlerFunc(srv.WSPublisher))
	defer s.Close()

	ws, err := connectWithWS(s)
	require.NoError(t, err)
	defer ws.Close()

	err = ws.WriteJSON(broker.MsgMessage("random_data"))
	assert.NoError(t, err)
}

func TestWSPublisherSyncNoSubscribers(t *testing.T) {
	srv := testServer()
	defer srv.Stop()

	pubHandler := httptest.NewServer(http.HandlerFunc(srv.WSPublisher))
	defer pubHandler.Close()

	publisher, err := connectWithWS(pubHandler)
	require.NoError(t, err)
	defer publisher.Close()

	err = publisher.WriteJSON(broker.MsgSyncSubscribers())
	assert.NoError(t, err)

	received := broker.Message{}
	err = publisher.ReadJSON(&received)
	assert.NoError(t, err)
	assert.Equal(t, broker.MsgNoSubscribers(), received)
}

func TestWSPublisherSyncHasSubscribers(t *testing.T) {
	srv := testServer()
	defer srv.Stop()

	pubHandler := httptest.NewServer(http.HandlerFunc(srv.WSPublisher))
	defer pubHandler.Close()

	subHandler := httptest.NewServer(http.HandlerFunc(srv.WSSubscriber))
	defer subHandler.Close()

	publisher, err := connectWithWS(pubHandler)
	require.NoError(t, err)
	defer publisher.Close()

	subscriber, err := connectWithWS(subHandler)
	require.NoError(t, err)
	defer subscriber.Close()

	err = publisher.WriteJSON(broker.MsgSyncSubscribers())
	assert.NoError(t, err)

	received := broker.Message{}
	err = publisher.ReadJSON(&received)
	assert.NoError(t, err)
	assert.Equal(t, broker.MsgHasSubscribers(), received)
}

func TestWSSubscriberReceives(t *testing.T) {
	srv := testServer()
	defer srv.Stop()

	pubHandler := httptest.NewServer(http.HandlerFunc(srv.WSPublisher))
	defer pubHandler.Close()

	subHandler := httptest.NewServer(http.HandlerFunc(srv.WSSubscriber))
	defer subHandler.Close()

	publisher, err := connectWithWS(pubHandler)
	assert.NoError(t, err)
	defer publisher.Close()

	subscriber, err := connectWithWS(subHandler)
	assert.NoError(t, err)
	defer subscriber.Close()

	publishedMsg := broker.MsgMessage("random_data")

	err = publisher.WriteJSON(publishedMsg)
	assert.NoError(t, err)

	receivedMsg := broker.Message{}
	err = subscriber.ReadJSON(&receivedMsg)
	assert.NoError(t, err)
	assert.Equal(t, publishedMsg, receivedMsg)
}

func TestMultipleWSSubscriberReceives(t *testing.T) {

	srv := testServer()
	defer srv.Stop()

	pubHandler := httptest.NewServer(http.HandlerFunc(srv.WSPublisher))
	defer pubHandler.Close()

	subHandler := httptest.NewServer(http.HandlerFunc(srv.WSSubscriber))
	defer subHandler.Close()

	publisher, err := connectWithWS(pubHandler)
	require.NoError(t, err)
	defer publisher.Close()

	subscriber1, err := connectWithWS(subHandler)
	require.NoError(t, err)
	defer subscriber1.Close()

	subscriber2, err := connectWithWS(subHandler)
	require.NoError(t, err)
	defer subscriber2.Close()

	publishedMsg := broker.MsgMessage("random_data")

	err = publisher.WriteJSON(publishedMsg)
	assert.NoError(t, err)

	receivedMsg := broker.Message{}

	err = subscriber1.ReadJSON(&receivedMsg)
	assert.NoError(t, err)
	assert.Equal(t, publishedMsg, receivedMsg)

	err = subscriber2.ReadJSON(&receivedMsg)
	assert.NoError(t, err)
	assert.Equal(t, publishedMsg, receivedMsg)
}

package wsreply

import "net/http"

func (a *Application) router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/sub", a.SubscriberWS)
	mux.HandleFunc("/pub", a.PublisherWS)

	return mux
}

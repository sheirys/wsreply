package wsreply

import (
	"net/http"
)

func (a *Application) router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/sub", a.WSSubscriber)
	mux.HandleFunc("/pub", a.WSPublisher)

	return mux
}

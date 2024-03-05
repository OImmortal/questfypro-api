package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Rota struct {
	Uri    string
	Metodo string
	Funcao func(w http.ResponseWriter, r *http.Request)
}

func Configurar(r *mux.Router) *mux.Router {
	routes := perguntas

	for _, rota := range routes {
		r.HandleFunc(rota.Uri, rota.Funcao).Methods(rota.Metodo)
	}

	return r
}

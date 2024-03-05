package routes

import (
	"api/src/controller"
	"net/http"
)

var perguntas = []Rota{
	{
		Uri:    "/perguntas",
		Metodo: http.MethodPost,
		Funcao: controller.CriarPergunta,
	},
	{
		Uri:    "/perguntas",
		Metodo: http.MethodGet,
		Funcao: func(w http.ResponseWriter, r *http.Request) {},
	},
	{
		Uri:    "/perguntas/{perguntaID}",
		Metodo: http.MethodGet,
		Funcao: func(w http.ResponseWriter, r *http.Request) {},
	},
}

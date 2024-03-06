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
		Funcao: controller.BuscarPergunta,
	},
	{
		Uri:    "/perguntas",
		Metodo: http.MethodGet,
		Funcao: controller.BuscarPerguntaCategoria,
	},
}

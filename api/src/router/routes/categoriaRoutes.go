package routes

import (
	"api/src/controller"
	"net/http"
)

var categorias = []Rota{
	{
		Uri:    "/categorias",
		Metodo: http.MethodPost,
		Funcao: controller.CriarCategoria,
	},
	{
		Uri:    "/categorias",
		Metodo: http.MethodGet,
		Funcao: controller.BuscarCategoria,
	},
}

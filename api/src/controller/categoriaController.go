package controller

import (
	"api/src/database"
	"api/src/models"
	"api/src/res"
	"encoding/json"
	"io"
	"net/http"
)

func CriarCategoria(w http.ResponseWriter, r *http.Request) {
	requestBody, erro := io.ReadAll(r.Body)
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	var categoria models.Categoria
	if erro = json.Unmarshal(requestBody, &categoria); erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	defer db.Close()

	statment, erro := db.Prepare("insert into categorias (categoria) value (?)")
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	defer statment.Close()
	_, erro = statment.Exec(categoria.Categoria)
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	res.JSON(w, http.StatusOK, "Categoria cadastrado com sucesso!!")
}

func BuscarCategoria(w http.ResponseWriter, r *http.Request) {

	db, erro := database.Conectar()
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	defer db.Close()

	linhas, erro := db.Query("select * from categorias")
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	defer linhas.Close()

	var categorias []models.Categoria
	for linhas.Next() {
		var categoria models.Categoria
		if erro = linhas.Scan(&categoria.Id, &categoria.Categoria); erro != nil {
			res.Erro(w, http.StatusBadRequest, erro)
			return
		}

		categorias = append(categorias, categoria)
	}

	res.JSON(w, http.StatusOK, categorias)
}

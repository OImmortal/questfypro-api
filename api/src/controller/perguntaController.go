package controller

import (
	"api/src/database"
	"api/src/models"
	"api/src/res"
	"encoding/json"
	"io"
	"net/http"
)

func CriarPergunta(w http.ResponseWriter, r *http.Request) {
	requestBody, erro := io.ReadAll(r.Body)
	if erro != nil {
		res.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	var pergunta models.Pergunta
	if erro = json.Unmarshal(requestBody, &pergunta); erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		res.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	statment, erro := db.Prepare("insert into perguntas (title, descrpt) values (?, ?)")
	if erro != nil {
		res.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer statment.Close()
	resultado, erro := statment.Exec(pergunta.Title, pergunta.Desc)
	if erro != nil {
		res.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	ultimoId, erro := resultado.LastInsertId()
	if erro != nil {
		res.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	for _, resposta := range pergunta.Resposta {
		statment, erro := db.Prepare("insert into respostas (id_pergunta, description, correta) values (?, ?, ?)")
		if erro != nil {
			res.Erro(w, http.StatusInternalServerError, erro)
			return
		}

		defer statment.Close()
		_, erro = statment.Exec(ultimoId, resposta.Desc, resposta.Correta)
		if erro != nil {
			res.Erro(w, http.StatusInternalServerError, erro)
			return
		}
	}

	res.JSON(w, http.StatusOK, "Pergunta cadastrada com sucesso")
}

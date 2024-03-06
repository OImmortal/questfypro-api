package controller

import (
	"api/src/database"
	"api/src/models"
	"api/src/res"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func CriarPergunta(w http.ResponseWriter, r *http.Request) {
	requestBody, erro := io.ReadAll(r.Body)
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	var pergunta models.Pergunta
	if erro = json.Unmarshal(requestBody, &pergunta); erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	defer db.Close()

	statment, erro := db.Prepare("insert into perguntas (title, descrpt, categoria_id) values (?, ?, ?)")
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	defer statment.Close()
	resultado, erro := statment.Exec(pergunta.Title, pergunta.Desc, pergunta.CategoriaId)
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	ultimoId, erro := resultado.LastInsertId()
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	for _, resposta := range pergunta.Resposta {
		statment, erro := db.Prepare("insert into respostas (id_pergunta, description, correta) values (?, ?, ?)")
		if erro != nil {
			res.Erro(w, http.StatusBadRequest, erro)
			return
		}

		defer statment.Close()
		_, erro = statment.Exec(ultimoId, resposta.Desc, resposta.Correta)
		if erro != nil {
			res.Erro(w, http.StatusBadRequest, erro)
			return
		}
	}

	res.JSON(w, http.StatusOK, "Pergunta cadastrada com sucesso")
}

func BuscarPergunta(w http.ResponseWriter, r *http.Request) {
	db, erro := database.Conectar()
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	defer db.Close()

	linhas, erro := db.Query("SELECT perguntas.id, title, descrpt, categoria_id, respostas.id, respostas.description, correta, categoria from perguntas inner join respostas on id_pergunta = perguntas.id INNER JOIN categorias ON categoria_id = categorias.id_categoria;")
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	defer linhas.Close()

	perguntas := make(map[int]models.Pergunta)

	for linhas.Next() {
		var perguntaID, perguntaCatId, respostaID int
		var perguntaTitle, perguntaDesc, respostaDesc, categoria string
		var respostaCorreta bool

		// Ler os valores da linha atual
		if erro := linhas.Scan(&perguntaID, &perguntaTitle, &perguntaDesc, &perguntaCatId, &respostaID, &respostaDesc, &respostaCorreta, &categoria); erro != nil {
			res.Erro(w, http.StatusBadRequest, erro)
			return
		}

		// Se a pergunta ainda não foi adicionada ao mapa, adicioná-la
		if _, ok := perguntas[perguntaID]; !ok {
			perguntas[perguntaID] = models.Pergunta{
				Id:            uint(perguntaID),
				Title:         perguntaTitle,
				Desc:          perguntaDesc,
				CategoriaNome: categoria,
				Resposta:      make([]models.Resposta, 0),
			}
		}

		// Adicionar a resposta à pergunta correspondente
		pergunta := perguntas[perguntaID]
		pergunta.Resposta = append(pergunta.Resposta, models.Resposta{
			Id:      uint(respostaID),
			Desc:    respostaDesc,
			Correta: respostaCorreta,
		})
		perguntas[perguntaID] = pergunta
	}

	// Converter o mapa em um slice de perguntas
	var perguntasSlice []models.Pergunta
	for _, pergunta := range perguntas {
		perguntasSlice = append(perguntasSlice, pergunta)
	}

	// Retornar as perguntas como JSON
	res.JSON(w, http.StatusOK, perguntasSlice)

}

func BuscarPerguntaCategoria(w http.ResponseWriter, r *http.Request) {
	categoria := strings.ToLower(r.URL.Query().Get("categoria"))

	db, erro := database.Conectar()
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	defer db.Close()

	linhas, erro := db.Query(
		"select id, title, descrpt, resposta.id, respostas.description, correta from perguntas inner join respostas on id_pergunta = perguntas.id where categoria = ?",
		categoria,
	)
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	defer linhas.Close()

	// perguntas := make(map[int]models.Pergunta)

	// for linhas.Next() {
	// 	var IdPergunta, IdResposta uint
	// 	var TitlePergunta, DescPergunta, DescResposta string
	// 	var CorretaResposta bool

	// }

}

package controller

import (
	"api/src/database"
	"api/src/models"
	"api/src/res"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
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

	categoria := strings.ToLower(r.URL.Query().Get("categoria"))
	if categoria == "" {
		perguntasSlice, erro := selectAll(w)
		if erro != nil {
			res.Erro(w, http.StatusBadRequest, erro)
		}
		res.JSON(w, http.StatusOK, perguntasSlice)
	} else {
		perguntasSlice, erro := selectByCategoria(w, categoria)
		if erro != nil {
			res.Erro(w, http.StatusBadRequest, erro)
		}
		res.JSON(w, http.StatusOK, perguntasSlice)
	}

}

func DeletePergunta(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	idPergunta, erro := strconv.ParseUint(parametros["idPergunta"], 10, 64)
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	defer db.Close()

	statment, erro := db.Prepare("delete from perguntas where id = ?")
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	defer statment.Close()
	if _, erro = statment.Exec(idPergunta); erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return
	}

	res.JSON(w, http.StatusOK, "Pergunta deletada com sucesso")
}

func selectAll(w http.ResponseWriter) ([]models.Pergunta, error) {
	db, erro := database.Conectar()
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return nil, erro
	}

	defer db.Close()

	linhas, erro := db.Query("SELECT perguntas.id, title, descrpt, categoria_id, respostas.id, respostas.description, correta, categoria from perguntas inner join respostas on id_pergunta = perguntas.id INNER JOIN categorias ON categoria_id = categorias.id_categoria;")
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return nil, erro
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
			return nil, erro
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

	return perguntasSlice, nil
}

func selectByCategoria(w http.ResponseWriter, categoria string) ([]models.Pergunta, error) {
	db, erro := database.Conectar()
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return nil, erro
	}

	defer db.Close()

	linhas, erro := db.Query(
		"SELECT perguntas.id, title, descrpt, respostas.id, respostas.description, correta from perguntas inner join respostas on id_pergunta = perguntas.id INNER JOIN categorias ON categoria_id = categorias.id_categoria WHERE categorias.categoria = ?",
		categoria,
	)
	if erro != nil {
		res.Erro(w, http.StatusBadRequest, erro)
		return nil, erro
	}

	defer linhas.Close()

	perguntas := make(map[int]models.Pergunta)

	for linhas.Next() {
		var IdPergunta, IdResposta int
		var TitlePergunta, DescPergunta, DescResposta string
		var CorretaResposta bool

		if erro = linhas.Scan(&IdPergunta, &TitlePergunta, &DescPergunta, &IdResposta, &DescResposta, &CorretaResposta); erro != nil {
			res.Erro(w, http.StatusBadRequest, erro)
			return nil, erro
		}

		if _, ok := perguntas[IdPergunta]; !ok {
			perguntas[IdPergunta] = models.Pergunta{
				Id:       uint(IdPergunta),
				Title:    TitlePergunta,
				Desc:     DescPergunta,
				Resposta: make([]models.Resposta, 0),
			}
		}

		pergunta := perguntas[IdPergunta]
		pergunta.Resposta = append(pergunta.Resposta, models.Resposta{
			Id:      uint(IdResposta),
			Desc:    DescResposta,
			Correta: CorretaResposta,
		})

		perguntas[IdPergunta] = pergunta
	}

	var perguntaSlice []models.Pergunta
	for _, pergunta := range perguntas {
		perguntaSlice = append(perguntaSlice, pergunta)
	}

	return perguntaSlice, nil
}

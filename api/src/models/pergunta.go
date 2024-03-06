package models

type Pergunta struct {
	Id            uint       `json:"id,omitempty"`
	Title         string     `json:"title"`
	Desc          string     `json:"description"`
	CategoriaId   uint       `json:"categoriaId,omitempty"`
	CategoriaNome string     `json:"categoriaNome,omitempty"`
	Resposta      []Resposta `json:"resposta"`
}

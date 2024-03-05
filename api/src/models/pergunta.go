package models

type Pergunta struct {
	Id       uint       `json:"id,omitempty"`
	Title    string     `json:"title"`
	Desc     string     `json:"description"`
	Resposta []Resposta `json:"resposta"`
}

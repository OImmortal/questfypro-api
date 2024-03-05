package models

type Resposta struct {
	Id      uint   `json:"id,omitempty"`
	Desc    string `json:"description"`
	Correta bool   `json:"correta"`
}

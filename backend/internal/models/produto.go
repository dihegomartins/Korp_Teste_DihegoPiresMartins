package models

// Produto representa a tabela de produto no SQL Server
type Produto struct {
	Id        int    `db:"Id" json:"id"`
	Codigo    string `db:"Codigo" json:"codigo"`
	Descricao string `db:"Descricao" json:"descricao"`
	Saldo     int    `db:"Saldo" json:"saldo"`
}
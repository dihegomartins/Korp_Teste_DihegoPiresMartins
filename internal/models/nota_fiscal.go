package models

// NotaFiscal reflete a tabela NotasFiscais
type NotaFiscal struct {
	Id               int              `json:"id" db:"Id"`
	NumeroSequencial int              `json:"numero_sequencial" db:"NumeroSequencial"`
	Status           string           `json:"status" db:"Status"`
	Itens            []NotaFiscalItem `json:"itens"` // Para receber os itens no POST
}

// NotaFiscalItem reflete a tabela NotaFiscalItens
type NotaFiscalItem struct {
	Id        int `json:"id" db:"Id"`
	NotaId    int `json:"nota_id" db:"NotaId"`
	ProdutoId int `json:"produto_id" db:"ProdutoId"`
	Quantidade int `json:"quantidade" db:"Quantidade"`
}
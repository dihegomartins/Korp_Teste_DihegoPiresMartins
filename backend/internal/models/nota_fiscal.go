package models

// NotaFiscal reflete a tabela NotasFiscais
type NotaFiscal struct {
    Id               int              `json:"Id" db:"Id"`
    NumeroSequencial int              `json:"NumeroSequencial" db:"NumeroSequencial"`
    Status           string           `json:"Status" db:"Status"` // Ajustado para S maiúsculo
    Itens            []NotaFiscalItem `json:"Itens" db:"-"`
}

// NotaFiscalItem reflete a tabela NotaFiscalItens
type NotaFiscalItem struct {
    Id               int    `json:"Id" db:"Id"`
    NotaId           int    `json:"NotaId" db:"NotaId"`
    ProdutoId        int    `json:"ProdutoId" db:"ProdutoId"`
    // O db:"-" diz ao Go para não tentar SALVAR isso na tabela NotaFiscalItens, 
    // já que a descrição vem da tabela de Produtos via JOIN
    DescricaoProduto string `json:"DescricaoProduto" db:"Descricao"` 
    Quantidade       int    `json:"Quantidade" db:"Quantidade"` // Ajustado para Q maiúsculo
}
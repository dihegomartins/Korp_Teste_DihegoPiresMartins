package repository

import (
	"fmt"
	"log"

	"github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/internal/database"
	"github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/internal/models"
)

func GetProdutoByCodigo(codigo string) (*models.Produto, error) {
	var produto models.Produto
	query := "SELECT Id, Codigo, Descricao, Saldo FROM Produtos WHERE Codigo = @p1"
	
	err := database.DB.Get(&produto, query, codigo)
	if err != nil {
		return nil, err
	}
	return &produto, nil
}

// Essencial para a baixa de estoque (ESCRITA)
func UpdateSaldo(id int, novoSaldo int) error {
	query := "UPDATE Produtos SET Saldo = @p1 WHERE Id = @p2"
	
	_, err := database.DB.Exec(query, novoSaldo, id)
	if err != nil {
		log.Printf("Erro ao atualizar saldo: %v", err)
		return err
	}
	return nil
}

// Útil para listagens (LEITURA)
func GetAllProdutos() ([]models.Produto, error) {
	var produtos []models.Produto
	query := "SELECT Id, Codigo, Descricao, Saldo FROM Produtos"
	
	err := database.DB.Select(&produtos, query)
	return produtos, err
}

// CreateProduto insere um novo produto no banco de dados
func CreateProduto(p models.Produto) error {
	query := `INSERT INTO Produtos (Codigo, Descricao, Saldo) 
              VALUES (@p1, @p2, @p3)`
	
	_, err := database.DB.Exec(query, p.Codigo, p.Descricao, p.Saldo)
	if err != nil {
		log.Printf("Erro ao inserir produto: %v", err)
		return err
	}
	return nil
}

// SubtrairSaldo reduz a quantidade do produto no banco
func SubtrairSaldo(id int, quantidadeParaSubtrair int) error {
	// A query subtrai do saldo atual. O "WHERE Saldo >= @p1" evita que o estoque fique negativo no banco
	query := `UPDATE Produtos SET Saldo = Saldo - @p1 WHERE Id = @p2 AND Saldo >= @p1`
	
	result, err := database.DB.Exec(query, quantidadeParaSubtrair, id)
	if err != nil {
		return err
	}

	// Verifica se alguma linha foi afetada (se o saldo era insuficiente, o SQL não atualiza nada)
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("saldo insuficiente ou produto não encontrado")
	}

	return nil
}
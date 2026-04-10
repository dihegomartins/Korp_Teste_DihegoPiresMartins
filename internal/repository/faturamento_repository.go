package repository

import (
	"context"
	"fmt"

	"github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/internal/models"
	"github.com/jmoiron/sqlx"
)

type FaturamentoRepository struct {
	db *sqlx.DB
}

func NewFaturamentoRepository(db *sqlx.DB) *FaturamentoRepository {
	return &FaturamentoRepository{db: db}
}

func (r *FaturamentoRepository) CriarNotaFiscalCompleta(nf models.NotaFiscal) error {
	ctx := context.Background()
	
	// 1. Inicia a Transação (Tratamento atômico: tudo ou nada)
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// 2. Insere a Nota Fiscal e recupera o ID gerado (Identity)
	var notaId int
	queryNota := `INSERT INTO NotasFiscais (NumeroSequencial, Status) 
				  OUTPUT INSERTED.Id 
				  VALUES (@p1, @p2)`
	
	err = tx.QueryRowxContext(ctx, queryNota, nf.NumeroSequencial, "Aberta").Scan(&notaId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao inserir nota: %v", err)
	}

	// 3. Itera sobre os itens da nota
	for _, item := range nf.Itens {
		// A. Baixa o estoque do produto apenas se houver saldo suficiente
		queryEstoque := `UPDATE Produtos SET Saldo = Saldo - @p1 
						 WHERE Id = @p2 AND Saldo >= @p1`
		
		res, err := tx.ExecContext(ctx, queryEstoque, item.Quantidade, item.ProdutoId)
		if err != nil {
			tx.Rollback()
			return err
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("produto ID %d sem saldo suficiente ou não encontrado", item.ProdutoId)
		}

		// B. Insere o item vinculado à nota criada no passo 2
		queryItem := `INSERT INTO NotaFiscalItens (NotaId, ProdutoId, Quantidade) 
					  VALUES (@p1, @p2, @p3)`
		_, err = tx.ExecContext(ctx, queryItem, notaId, item.ProdutoId, item.Quantidade)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 4. Finaliza a transação salvando tudo definitivamente
	return tx.Commit()
}
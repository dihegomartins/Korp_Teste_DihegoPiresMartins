package repository

import (
	"context"
	"database/sql"
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
    
    tx, err := r.db.BeginTxx(ctx, nil)
    if err != nil {
        return err
    }

    var notaId int
    queryNota := `INSERT INTO NotasFiscais (NumeroSequencial, Status) 
                  OUTPUT INSERTED.Id 
                  VALUES (@p1, @p2)`
    
    err = tx.QueryRowxContext(ctx, queryNota, nf.NumeroSequencial, "Aberta").Scan(&notaId)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("erro ao inserir nota: %v", err)
    }

    for _, item := range nf.Itens {

        queryItem := `INSERT INTO NotaFiscalItens (NotaId, ProdutoId, Quantidade) 
                      VALUES (@p1, @p2, @p3)`
        _, err = tx.ExecContext(ctx, queryItem, notaId, item.ProdutoId, item.Quantidade)
        if err != nil {
            tx.Rollback()
            return err
        }
    }
    return tx.Commit()
}

func (r *FaturamentoRepository) ListarNotas() ([]models.NotaFiscal, error) {
	var notas []models.NotaFiscal
	
	queryNotas := `SELECT Id, NumeroSequencial, Status FROM NotasFiscais`
	err := r.db.Select(&notas, queryNotas)
	if err != nil {
		return nil, err
	}

	for i := range notas {
		var itens []models.NotaFiscalItem
		
		// 1. Escrevemos a query com ?
		rawQuery := `
			SELECT 
				ni.Id, 
				ni.NotaId, 
				ni.ProdutoId, 
				p.Descricao, 
				ni.Quantidade 
			FROM NotaFiscalItens ni
			INNER JOIN Produtos p ON ni.ProdutoId = p.Id
			WHERE ni.NotaId = ?`
		
		// 2. O Rebind transforma o ? no @p1 que o SQL Server exige
		queryItens := r.db.Rebind(rawQuery)

		// 3. Executamos
		err := r.db.Select(&itens, queryItens, notas[i].Id)
		if err != nil {
			return nil, err
		}
		notas[i].Itens = itens
	}

	return notas, nil
}


func (r *FaturamentoRepository) FecharNotaFiscal(id int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	// 1. Verificar se a nota existe e está "Aberta" 
	var status string
	err = tx.Get(&status, "SELECT Status FROM NotasFiscais WHERE Id = @p1", id)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return fmt.Errorf("nota fiscal com ID %d não encontrada", id)
		}
		return err
	}
	if status != "Aberta" {
		tx.Rollback()
		return fmt.Errorf("apenas notas com status 'Aberta' podem ser fechadas")
	}

	// 2. Buscar itens da nota para baixar o estoque [cite: 30, 37]
	var itens []models.NotaFiscalItem
	err = tx.Select(&itens, "SELECT ProdutoId, Quantidade FROM NotaFiscalItens WHERE NotaId = @p1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 3. Atualizar saldo de cada produto [cite: 37, 38]
	for _, item := range itens {
		res, err := tx.Exec("UPDATE Produtos SET Saldo = Saldo - @p1 WHERE Id = @p2 AND Saldo >= @p1", item.Quantidade, item.ProdutoId)
		if err != nil {
			tx.Rollback()
			return err
		}
		rows, _ := res.RowsAffected()
		if rows == 0 {
			tx.Rollback()
			return fmt.Errorf("produto ID %d sem saldo suficiente para fechar a nota", item.ProdutoId)
		}
	}

	// 4. Atualizar status da nota para "Fechada" 
	_, err = tx.Exec("UPDATE NotasFiscais SET Status = 'Fechada' WHERE Id = @p1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
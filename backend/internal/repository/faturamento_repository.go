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
	// 1. Insere a Nota (Garante que NumeroSequencial e Status batam com o script)
	queryNota := r.db.Rebind(`INSERT INTO NotasFiscais (NumeroSequencial, Status) 
				  OUTPUT INSERTED.Id 
				  VALUES (?, ?)`)
	
	err = tx.QueryRowxContext(ctx, queryNota, nf.NumeroSequencial, "Aberta").Scan(&notaId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao inserir nota: %v", err)
	}

	// 2. Insere os Itens
	for _, item := range nf.Itens {
		// Verifique se no seu banco está NotaFiscalItens ou NotaFiscallItens
		queryItem := r.db.Rebind(`INSERT INTO NotaFiscalItens (NotaId, ProdutoId, Quantidade) 
					  VALUES (?, ?, ?)`)
		
		_, err = tx.ExecContext(ctx, queryItem, notaId, item.ProdutoId, item.Quantidade)
		if err != nil {
			tx.Rollback()
			// Este log vai te mostrar no terminal do VS Code qual ProdutoId deu erro
			return fmt.Errorf("erro no item (ProdutoId: %d): %v", item.ProdutoId, err)
		}
	}
	
	return tx.Commit()
}

func (r *FaturamentoRepository) ListarNotas() ([]models.NotaFiscal, error) {
	var notas []models.NotaFiscal
	
	// Ajustado NumeroSequencial
	queryNotas := `SELECT Id, NumeroSequencial, Status FROM NotasFiscais`
	err := r.db.Select(&notas, queryNotas)
	if err != nil {
		return nil, err
	}

	for i := range notas {
		var itens []models.NotaFiscalItem
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
		
		err := r.db.Select(&itens, r.db.Rebind(rawQuery), notas[i].Id)
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

	var status string
	err = tx.Get(&status, r.db.Rebind("SELECT Status FROM NotasFiscais WHERE Id = ?"), id)
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

	_, err = tx.Exec(r.db.Rebind("UPDATE NotasFiscais SET Status = 'Fechada' WHERE Id = ?"), id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *FaturamentoRepository) BuscarNotaPorID(id int) (*models.NotaFiscal, error) {
	var nota models.NotaFiscal
	queryNota := r.db.Rebind(`SELECT Id, NumeroSequencial, Status FROM NotasFiscais WHERE Id = ?`)
	err := r.db.Get(&nota, queryNota, id)
	if err != nil {
		return nil, err
	}

	queryItens := r.db.Rebind(`SELECT Id, NotaId, ProdutoId, Quantidade FROM NotaFiscalItens WHERE NotaId = ?`)
	err = r.db.Select(&nota.Itens, queryItens, id)
	if err != nil {
		return &nota, nil
	}

	return &nota, nil
}
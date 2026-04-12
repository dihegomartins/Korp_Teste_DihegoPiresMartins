package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/internal/database"
	"github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/internal/models"
	"github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/internal/repository"
	"github.com/gin-gonic/gin"
)

// AbrirNotaFiscal com validações de campos obrigatórios
func AbrirNotaFiscal(c *gin.Context) {
	var nf models.NotaFiscal
	if err := c.ShouldBindJSON(&nf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de dados inválido: " + err.Error()})
		return
	}

	// [Validação Obrigatória] Verifica numeração sequencial
	if nf.NumeroSequencial <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A Numeração Sequencial é obrigatória e deve ser maior que zero"})
		return
	}

	// [Regra de Negócio] Garante status inicial como "Aberta"
	nf.Status = "Aberta"

	db := database.GetDB()
	repo := repository.NewFaturamentoRepository(db)
	
	err := repo.CriarNotaFiscalCompleta(nf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao persistir nota fiscal: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Nota Fiscal nº %d criada com sucesso (Status: Aberta)!", nf.NumeroSequencial),
	})
}

func ListarNotas(c *gin.Context) {
	db := database.GetDB()
	repo := repository.NewFaturamentoRepository(db)
	notas, err := repo.ListarNotas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar notas: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, notas)
}

func FecharNota(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	db := database.GetDB()
	repo := repository.NewFaturamentoRepository(db)

	nota, err := repo.BuscarNotaPorID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nota fiscal não encontrada"})
		return
	}

	if nota.Status != "Aberta" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Apenas notas com status Aberta podem ser fechadas/impressas"})
		return
	}

	client := &http.Client{}
	for _, item := range nota.Itens {
		url := fmt.Sprintf("http://127.0.0.1:8081/produtos/%d/baixa", item.ProdutoId)
		
		req, err := http.NewRequest(http.MethodPatch, url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao preparar comunicação com estoque"})
			return
		}

		req.Header.Set("X-Quantidade-Baixa", strconv.Itoa(item.Quantidade))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)

		if err != nil {
			fmt.Printf("❌ Falha de conexão: %v\n", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Falha na integração: O microsserviço de estoque está offline."})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var bodyErr struct {
				Error string `json:"error"`
			}
			json.NewDecoder(resp.Body).Decode(&bodyErr)
			
			mensagemFinal := bodyErr.Error
			if mensagemFinal == "" {
				mensagemFinal = "Erro desconhecido no microsserviço de estoque."
			}

			c.JSON(resp.StatusCode, gin.H{
				"error": "Integração Rejeitada: " + mensagemFinal,
			})
			return
		}
	}

	err = repo.FecharNotaFiscal(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao atualizar status final: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nota finalizada e estoque atualizado com sucesso!"})
}
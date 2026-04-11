package handlers

import (
	"net/http"
	"strconv"

	"github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/internal/database"
	"github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/internal/models"
	"github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/internal/repository"
	"github.com/gin-gonic/gin"
)

// AbrirNotaFiscal godoc
// @Summary      Abrir nova Nota Fiscal
// @Description  Cria uma nota fiscal, abate o estoque dos produtos e salva os itens em uma transação única.
// @Tags         faturamento
// @Accept       json
// @Produce      json
// @Param        nota  body      models.NotaFiscal  true  "Dados da Nota e Itens"
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /faturamento [post]
func AbrirNotaFiscal(c *gin.Context) {
	var nf models.NotaFiscal

	// 1. Faz o Bind do JSON enviado pelo Angular/Insonmia
	if err := c.ShouldBindJSON(&nf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados da nota inválidos: " + err.Error()})
		return
	}

	// 2. Instancia o repositório (usando a conexão global que você já tem no database.go)
	db := database.GetDB()
	repo := repository.NewFaturamentoRepository(db)

	// 3. Chama a função de transação que criamos no repositório
	err := repo.CriarNotaFiscalCompleta(nf)
	if err != nil {
		// Se der erro (ex: falta de estoque), a transação já deu Rollback lá no repo
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Retorno de sucesso
	c.JSON(http.StatusCreated, gin.H{
		"message": "Nota Fiscal nº " + string(rune(nf.NumeroSequencial)) + " aberta com sucesso e estoque atualizado!",
	})
}

// ListarNotas godoc
// @Summary      Listar todas as notas
// @Description  Retorna o histórico de notas fiscais com seus itens
// @Tags         faturamento
// @Produce      json
// @Success      200  {array}  models.NotaFiscal
// @Router       /faturamento [get]
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


// FecharNota godoc
// @Summary      Fechar e Imprimir Nota Fiscal
// @Description  Atualiza o status para Fechada e deduz o estoque dos produtos 
// @Tags         faturamento
// @Param        id   path      int  true  "ID da Nota Fiscal"
// @Success      200  {object}  map[string]string
// @Router       /faturamento/{id}/fechar [patch]
func FecharNota(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	db := database.GetDB()
	repo := repository.NewFaturamentoRepository(db)

	err := repo.FecharNotaFiscal(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nota finalizada e estoque atualizado com sucesso!"})
}
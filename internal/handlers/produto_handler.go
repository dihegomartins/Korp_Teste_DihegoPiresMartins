package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/internal/models"
	"github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/internal/repository"
	"github.com/gin-gonic/gin"
)

// GetProduto recebe a requisição e chama o repositório
func GetProduto(c *gin.Context) {
	// Pega o código que virá na URL: /produtos/:codigo
	codigo := c.Param("codigo")

	produto, err := repository.GetProdutoByCodigo(codigo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produto não encontrado"})
		return
	}

	// Se deu tudo certo, devolve o produto como JSON
	c.JSON(http.StatusOK, produto)
}

// UpdateEstoque recebe um JSON com o novo saldo e atualiza no banco
func UpdateEstoque(c *gin.Context) {
	// 1. Pegamos o ID da URL
	idStr := c.Param("id")
	
	// Convertendo o ID para inteiro
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	// 2. Pegamos o novo saldo do corpo da requisição (JSON)
	var body struct {
		NovoSaldo int `json:"novo_saldo"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	// 3. Chamamos o repositório para escrever no banco
	err := repository.UpdateSaldo(id, body.NovoSaldo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar banco"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estoque atualizado com sucesso!"})
}

// ListarProdutos devolve todos os produtos para o front-end
func ListarProdutos(c *gin.Context) {
	produtos, err := repository.GetAllProdutos()
	if err != nil {
		// Log interno para o desenvolvedor ver no terminal
		log.Printf("ERRO CRÍTICO DE BANCO: %v", err)

		// Resposta clara para o cliente (Angular)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "O serviço de banco de dados está temporariamente indisponível. Tente novamente em instantes.",
		})
		return
	}

	// Retorna a lista (se estiver vazia, retorna [])
	c.JSON(http.StatusOK, produtos)
}

// CriarProduto lida com a requisição POST para novos produtos
func CriarProduto(c *gin.Context) {
	var novoProduto models.Produto

	// 1. Tenta "bindar" o JSON recebido na nossa struct
	if err := c.ShouldBindJSON(&novoProduto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// 2. Chama o repositório para salvar
	err := repository.CreateProduto(novoProduto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível salvar o produto"})
		return
	}

	// 3. Retorna sucesso
	c.JSON(http.StatusCreated, gin.H{"message": "Produto cadastrado com sucesso!", "data": novoProduto})
}

// BaixaEstoque lida com a retirada de itens do inventário
func BaixaEstoque(c *gin.Context) {
	idStr := c.Param("id")
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	var request struct {
		Quantidade int `json:"quantidade"`
	}

	if err := c.ShouldBindJSON(&request); err != nil || request.Quantidade <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantidade inválida para baixa"})
		return
	}

	// Chama o repositório para subtrair
	err := repository.SubtrairSaldo(id, request.Quantidade)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Baixa realizada com sucesso!"})
}
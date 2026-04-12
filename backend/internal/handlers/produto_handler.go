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

// ListarProdutos godoc
// @Summary      Listar produtos
// @Description  Retorna todos os produtos do banco de dados
// @Tags         produtos
// @Produce      json
// @Success      200  {array}  models.Produto
// @Router       /produtos [get]
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

// CriarProduto godoc
// @Summary      Cadastrar novo produto
// @Description  Salva um novo produto no banco de dados
// @Tags         produtos
// @Accept       json
// @Produce      json
// @Param        produto  body      models.Produto  true  "Dados do Produto"
// @Success      201      {object}  models.Produto
// @Router       /produtos [post]
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

// BaixaEstoque godoc
// @Summary      Dar baixa no estoque
// @Description  Subtrai uma quantidade do saldo (aceita via JSON ou Header para integração)
// @Tags         produtos
// @Accept       json
// @Produce      json
// @Param        id           path      int  true  "ID do Produto"
// @Success      200          {object}  map[string]string
// @Router       /produtos/{id}/baixa [patch]
func BaixaEstoque(c *gin.Context) {
	idStr := c.Param("id")
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	var quantidade int

	// 1. Tenta pegar a quantidade do HEADER (Enviada pelo Microsserviço de Faturamento)
	qtdHeader := c.GetHeader("X-Quantidade-Baixa")
	if qtdHeader != "" {
		fmt.Sscanf(qtdHeader, "%d", &quantidade)
	}

	// 2. Se não estiver no Header, tenta pegar do JSON (Enviada pelo Angular/Insomnia)
	if quantidade <= 0 {
		var request struct {
			Quantidade int `json:"quantidade"`
		}
		// BindJSON pode falhar se o corpo estiver vazio, por isso ignoramos o erro aqui 
		// e validamos a variável 'quantidade' logo abaixo
		_ = c.ShouldBindJSON(&request)
		if request.Quantidade > 0 {
			quantidade = request.Quantidade
		}
	}

	// 3. Validação final
	if quantidade <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantidade inválida para baixa (deve ser maior que zero)"})
		return
	}

	// 4. Chama o repositório para subtrair
	err := repository.SubtrairSaldo(id, quantidade)
	if err != nil {
		// Se o repositório retornar erro (ex: saldo insuficiente), enviamos 409 Conflict
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Baixa realizada com sucesso!"})
}


// AdicionarEstoque godoc
// @Summary      Reposição de estoque
// @Description  Soma uma quantidade ao saldo atual de um produto específico
// @Tags         produtos
// @Accept       json
// @Produce      json
// @Param        id          path      int  true  "ID do Produto"
// @Param        quantidade  body      object  true  "Quantidade a adicionar (ex: {"quantidade": 50})"
// @Success      200         {object}  map[string]string
// @Router       /produtos/{id}/adicionar [patch]
func AdicionarEstoque(c *gin.Context) {
	// 1. Pegamos o ID da URL
	idStr := c.Param("id")
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	// 2. Pegamos a quantidade do corpo da requisição
	var request struct {
		Quantidade int `json:"quantidade"`
	}

	if err := c.ShouldBindJSON(&request); err != nil || request.Quantidade <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantidade inválida para reposição"})
		return
	}

	// 3. Chama o repositório que acabamos de atualizar para SOMAR ao saldo
	err := repository.AdicionarSaldo(id, request.Quantidade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível adicionar o saldo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estoque reposto com sucesso!"})
}
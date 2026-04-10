package main

import (
	"log"

	_ "github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/docs"
	"github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/internal/database"
	"github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// 1. Carrega o arquivo .env
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	// 2. Conecta ao Banco
	database.Connect()

	// 3. Inicializa o Servidor Gin
	router := gin.Default()
	router.Use(cors.Default())


	// Rota de teste: Ver se a API está viva
	// Rota da documentação
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API de Estoque está online! 🚀",
		})
	})
	router.GET("/produtos/:codigo", handlers.GetProduto)
	router.GET("/produtos", handlers.ListarProdutos)
	router.PUT("/produtos/:id/estoque", handlers.UpdateEstoque)
	router.POST("/produtos", handlers.CriarProduto)
	router.PATCH("/produtos/:id/baixa", handlers.BaixaEstoque)
	
	// 4. Sobe o servidor na porta 8080
	log.Println("🚀 Servidor rodando em http://localhost:8080")
	router.Run(":8080")
}
package main

import (
	"log"
	"os"

	_ "github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/docs"
	"github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/internal/database"
	"github.com/dihegomartins/Korp_Teste_DihegoPiresMartins/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// 1. Tenta carregar o .env de múltiplos locais possíveis
	// Primeiro tenta subir dois níveis (se rodado de dentro de cmd/estoque)
	err := godotenv.Load("../../.env")
	if err != nil {
		// Se falhar, tenta na raiz atual (se rodado da pasta /backend)
		err = godotenv.Load(".env")
		if err != nil {
			log.Println("Aviso: Arquivo .env não encontrado. Verificando variáveis de ambiente do sistema...")
		}
	}

	// Verificação Crítica: Se a URL do banco não existir, o serviço falha com feedback
	if os.Getenv("DB_URL") == "" {
		log.Fatal("❌ Erro Crítico: Variável de ambiente DB_URL não definida! Verifique seu arquivo .env")
	}

	// 2. Conecta ao Banco de Dados Real
	database.Connect()

	// 3. Inicializa o Servidor Gin
	router := gin.Default()
	
	// Middleware de CORS para comunicação com o Frontend Angular
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	produtos := router.Group("/produtos")
	{
		produtos.GET("", handlers.ListarProdutos) 
		produtos.POST("", handlers.CriarProduto) 
		produtos.GET("/:codigo", handlers.GetProduto)

		produtos.PATCH("/:id/adicionar", handlers.AdicionarEstoque) 
		router.PATCH("/produtos/:id/baixa", handlers.BaixaEstoque)
	}

	// 4. Sobe o serviço na porta 8081
	log.Println("📦 Microsserviço de ESTOQUE rodando em http://localhost:8081")
	log.Fatal(router.Run(":8081"))
}
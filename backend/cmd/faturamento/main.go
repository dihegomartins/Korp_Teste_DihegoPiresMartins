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
	// 1. Carrega o .env de múltiplos locais possíveis para evitar erros de path
	err := godotenv.Load("../../.env") 
	if err != nil {
		err = godotenv.Load(".env") 
		if err != nil {
			log.Println("Aviso: Arquivo .env não encontrado. Verificando variáveis de ambiente do sistema...")
		}
	}

	// Validação Crítica: Garante conexão real com o banco 
	if os.Getenv("DB_URL") == "" {
		log.Fatal("❌ Erro Crítico: Variável de ambiente DB_URL não definida!")
	}

	// 2. Conecta ao Banco de Dados (Persistência Física) 
	database.Connect()

	// 3. Inicializa o Servidor Gin 
	router := gin.Default()
	
	// Middleware de CORS para integração com Angular
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
	faturamento := router.Group("/faturamento")
	{
		faturamento.GET("", handlers.ListarNotas)
		faturamento.POST("", handlers.AbrirNotaFiscal)
		faturamento.PATCH("/:id/fechar", handlers.FecharNota)
	}

	// 4. Sobe o serviço na porta 8082 
	log.Println("🧾 Microsserviço de FATURAMENTO rodando em http://localhost:8082")
	log.Fatal(router.Run(":8082"))
}
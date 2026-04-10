package database

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"
)

var DB *sqlx.DB

func Connect() {
	var err error
	// Pegamos a URL do .env
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("A variável de ambiente DB_URL não foi definida no arquivo .env")
	}

	// Abrimos a conexão
	DB, err = sqlx.Connect("sqlserver", dsn)
	if err != nil {
		log.Fatal("Erro ao configurar o driver do banco:", err)
	}

	// O 'Ping' é o que realmente testa se os dados de sa:senha@localhost estão certos
	err = DB.Ping()
	if err != nil {
		log.Fatal("Não foi possível conectar ao banco de dados. Verifique se o SQL Server está rodando e se a senha está correta. Erro:", err)
	}

	log.Println("✅ Conexão com SQL Server (KorpDB) estabelecida com sucesso!")
}

// GetDB retorna a instância de conexão com o banco de dados
func GetDB() *sqlx.DB {
    return DB
}
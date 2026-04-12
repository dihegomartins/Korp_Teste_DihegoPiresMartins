# Desafio Korp - Sistema de Faturamento e Estoque

Este projeto é uma solução de microserviços para controle de estoque e faturamento, desenvolvido como parte do desafio técnico para a Korp. O sistema permite a criação de notas fiscais, gerenciamento de itens e atualização automática de estoque através de comunicação entre serviços.

## 🚀 Tecnologias Utilizadas

### Backend

- **Go (Golang)**: Linguagem principal para os microserviços.
- **SQL Server**: Banco de dados relacional.
- **Sqlx**: Extensão para interação com o banco de dados (mapeamento de structs).
- **Context API**: Gerenciamento de timeouts e controle de concorrência.

### Frontend

- **Angular 17+**: Framework para a interface do usuário (Standalone Components).
- **Bootstrap / Bootstrap Icons**: Estilização e componentes visuais.
- **RxJS**: Programação reativa para chamadas de API.

---

## 📦 Estrutura do Projeto

- `/backend`: Contém o código-fonte Go, organizado em `cmd` para os pontos de entrada.
  - `cmd/estoque`: Microserviço de Estoque (Porta **8081**).
  - `cmd/faturamento`: Microserviço de Faturamento (Porta **8082**).
- `/frontend`: Aplicação Angular (Porta **4200**).

---

## 🛠️ Como Executar o Projeto

### 1. Banco de Dados

Execute o script SQL fornecido para criar o banco `KorpDB` e as tabelas:

- `Produtos`
- `NotasFiscais`
- `NotaFiscalItens`

### 2. Configuração do Backend

1. Navegue até a pasta raiz `/backend`.
2. Configure o arquivo `.env` com a sua string de conexão do SQL Server:
   `DB_URL="sqlserver://usuario:senha@localhost:1433?database=KorpDB"`
3. Para rodar os serviços, abra dois terminais na pasta `/backend` e execute:
   - **Serviço de Estoque:** `go run cmd/estoque/main.go`
   - **Serviço de Faturamento:** `go run cmd/faturamento/main.go`

### 3. Configuração do Frontend

1. Navegue até a pasta `/frontend`.
2. Instale as dependências: `npm install`
3. Inicie a aplicação: `ng serve`
4. Acesse no navegador: `http://localhost:4200`

---

## 📝 Requisitos Atendidos

- [x] **Criação de Nota Fiscal**: Permite informar numeração sequencial personalizada.
- [x] **Status Inicial**: Toda nota é criada automaticamente com o status "Aberta".
- [x] **Múltiplos Itens**: Interface dinâmica para adicionar vários produtos em uma única nota antes de salvar.
- [x] **Integração de Microserviços**: Ao "Imprimir" (fechar) a nota, o serviço de faturamento comunica-se com o de estoque para realizar a baixa dos saldos.
- [x] **Listagem em Tempo Real**: Visualização de notas e consulta de saldo atualizado de produtos.

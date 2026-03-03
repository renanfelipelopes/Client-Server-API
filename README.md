# 💵 Desafio Go --- CLIENT-SERVER-API

Este projeto implementa:

-   🖥 **Servidor HTTP** que:
    -   Consome API externa de cotação
    -   Aplica timeout de 200ms na API
    -   Persiste a cotação no banco SQLite com timeout de 10ms
    -   Retorna apenas o campo `bid`
-   🧑‍💻 **Client HTTP** que:
    -   Possui timeout de 300ms
    -   Consome o servidor
    -   Salva a cotação em arquivo `cotacao.txt`

------------------------------------------------------------------------

# 📌 Arquitetura

Client (300ms timeout) ↓ Server (/cotacao) ↓ API Externa (200ms timeout)
↓ SQLite (10ms timeout)

Cada camada possui seu próprio controle de timeout, garantindo
isolamento e resiliência.

------------------------------------------------------------------------

# 🚀 Tecnologias Utilizadas

-   Go 1.22+
-   net/http
-   context
-   database/sql
-   SQLite
-   Driver: `modernc.org/sqlite` (sem CGO)

------------------------------------------------------------------------

# 📂 Estrutura do Projeto

. ├── server.go ├── client.go ├── go.mod └── README.md

------------------------------------------------------------------------

# ⚙️ Como Rodar o Projeto Localmente

## 1️⃣ Clonar o projeto

git clone https://github.com/renanfelipelopes/Client-Server-API

------------------------------------------------------------------------

## 2️⃣ Instalar dependências

go mod tidy

------------------------------------------------------------------------

## 3️⃣ Rodar o servidor

go run server.go

Servidor iniciará em:

http://localhost:8080

Endpoint disponível:

GET /cotacao

------------------------------------------------------------------------

## 4️⃣ Rodar o client (em outro terminal)

go run client.go

------------------------------------------------------------------------

# 📄 Resultado Esperado

Após executar o client, será criado o arquivo:

cotacao.txt

Conteúdo obrigatório:

Dólar: 5.1234

(O valor varia conforme a cotação atual)

------------------------------------------------------------------------

# ⏱️ Timeouts Implementados

## 🔹 API Externa (Servidor)

context.WithTimeout(r.Context(), 200\*time.Millisecond)

Se exceder: - Retorna HTTP 504 - Loga erro no console

------------------------------------------------------------------------

## 🔹 Banco de Dados (SQLite)

context.WithTimeout(context.Background(), 10\*time.Millisecond)

Se exceder: - Loga erro no console

------------------------------------------------------------------------

## 🔹 Client

context.WithTimeout(context.Background(), 300\*time.Millisecond)

Se exceder: - Loga `context deadline exceeded`

------------------------------------------------------------------------

# 🗄 Banco de Dados

Banco utilizado: **SQLite**

Arquivo gerado automaticamente:

cotacao.db

Tabela criada automaticamente:

CREATE TABLE IF NOT EXISTS cotacoes ( id INTEGER PRIMARY KEY
AUTOINCREMENT, bid TEXT, created_at DATETIME );

------------------------------------------------------------------------

# ⚠️ Sobre CGO e SQLite

Este projeto utiliza o driver:

modernc.org/sqlite

Motivo:

-   Não depende de CGO
-   Funciona em Docker Alpine
-   Funciona com CGO_ENABLED=0
-   Mais simples para deploy

------------------------------------------------------------------------

# 🧪 Teste Manual do Endpoint

Com o servidor rodando:

curl http://localhost:8080/cotacao

Resposta esperada:

{ "bid": "5.1234" }

------------------------------------------------------------------------

# 🎯 Requisitos Atendidos

✔ Timeout de 200ms na API externa\
✔ Timeout de 10ms no banco\
✔ Timeout de 300ms no client\
✔ Persistência no SQLite\
✔ Criação automática da tabela\
✔ Salvamento em arquivo `cotacao.txt`\
✔ Tratamento correto de erros\
✔ Sem uso de panic\
✔ Logs apropriados

------------------------------------------------------------------------

# 🧠 Conceitos Demonstrados

-   Context propagation
-   Timeout control
-   HTTP client/server
-   JSON encoding/decoding
-   ExecContext no banco
-   Separação de responsabilidades
-   Tratamento robusto de erro

------------------------------------------------------------------------

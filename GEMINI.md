# Resumo do Projeto: Zion API

Este documento fornece um resumo abrangente da Zion API, uma aplicação Go (Golang) que segue os princípios da Clean Architecture. O objetivo é servir como uma base de conhecimento para facilitar o desenvolvimento e a manutenção contínua.

## 1. Visão Geral e Tecnologias

-   **Linguagem:** Go (versão 1.23.3)
-   **Framework Web:** Gin (v1.10.0)
-   **Banco de Dados:** PostgreSQL
-   **Driver de Banco de Dados:** pgx (v5.7.1)
-   **Query Builder:** Squirrel (v1.5.4) para construção de queries SQL seguras e fluentes.
-   **Autenticação:** JWT com a plataforma Clerk, validando tokens através de uma chave pública PEM.
-   **Variáveis de Ambiente:** `godotenv` para carregar variáveis de `.env`.
-   **Tratamento de CORS:** `gin-contrib/cors` para permitir requisições de origens configuradas.
-   **Ponto de Entrada:** `cmd/api/main.go`

## 2. Arquitetura

O projeto adota a **Clean Architecture**, dividindo as responsabilidades em camadas claras:

-   **`internal/domain`**: Contém as entidades de negócio (structs Go) e as interfaces dos repositórios. É o coração da aplicação, sem dependências externas.
-   **`internal/usecase`**: Implementa a lógica de negócio pura, orquestrando as operações sobre as entidades. Depende apenas das interfaces do domínio.
-   **`internal/controller`**: Responsável por lidar com as requisições HTTP (camada de apresentação). Recebe dados do Gin, chama os `usecases` e retorna as respostas.
-   **`internal/infra`**: Implementa as dependências externas, como os repositórios de banco de dados.
    -   **`repository/postgres`**: Implementação concreta das interfaces de repositório do domínio, utilizando PostgreSQL e Squirrel.
    -   **`database`**: Gerencia a conexão com o banco de dados.
-   **`internal/middleware`**: Contém os middlewares HTTP, como o de autenticação (`AuthMiddleware`).
-   **`cmd/api/main.go`**: Ponto de entrada da aplicação. Configura o Gin, o pool de conexões com o banco, o CORS, as rotas e injeta as dependências (DI - Dependency Injection).

## 3. Modelos de Dados e Esquema do Banco

As migrações SQL revelam as seguintes tabelas e seus relacionamentos:

### Tabela: `products`

-   **Campos:** `id` (UUID), `name` (varchar), `value` (int), `unity_type` (char(2)), `category_id` (UUID), `is_deleted` (bool), `created_at`, `updated_at`.
-   **Relacionamentos:** Possui uma chave estrangeira `category_id` que referencia `category_products(id)`.

### Tabela: `customers`

-   **Campos:** `id` (UUID), `name` (text), `phone` (text), `phone2` (text), `email` (text), `is_deleted` (bool), `created_at`, `updated_at`.

### Tabela: `addresses`

-   **Campos:** `id` (UUID), `customer_id` (UUID), `cep` (text), `street` (text), `number` (text), `neighborhood` (text), `city` (text), `state` (text), `aditional_details` (text), `distance` (int), `is_default` (bool), `created_at`, `updated_at`.
-   **Relacionamentos:** Possui uma chave estrangeira `customer_id` que referencia `customers(id)`.

### Tabela: `category_products`

-   **Campos:** `id` (UUID), `name` (text), `description` (text).
-   **Dados Iniciais:** A tabela é populada com categorias padrão como 'Diversos', 'Massas', 'Bebidas', etc.

## 4. Fluxo de uma Requisição (Exemplo: `GET /products`)

1.  **`main.go`**: A rota `/products` é registrada e associada ao `productController.GetAll`. A rota está dentro de um grupo protegido pelo `AuthMiddleware`.
2.  **`middleware/auth.go`**: O `AuthMiddleware` intercepta a requisição, extrai o token JWT do cabeçalho `Authorization`, e o valida usando a chave pública do Clerk. Se inválido, a requisição é bloqueada.
3.  **`controller/product.go`**: O método `GetAll` do `productController` é invocado. Ele extrai os parâmetros de paginação da query string da URL.
4.  **`usecase/product.go`**: O controller chama o `GetAll` do `productUseCase`, passando os dados de paginação.
5.  **`usecase/product.go`**: O `productUseCase` chama o método `FindAll` da interface `ProductRepository` (injetada em sua construção).
6.  **`infra/repository/postgres/pg_product_repository.go`**: A implementação concreta do repositório é executada. Ela usa o **Squirrel** para construir a query `SELECT` para a tabela `products`, incluindo a lógica de paginação (`LIMIT`, `OFFSET`).
7.  **Banco de Dados**: A query é executada no PostgreSQL.
8.  **Retorno**: Os dados retornam pela mesma cadeia de chamadas, sendo transformados de structs de banco para DTOs (Data Transfer Objects) ou entidades de domínio, e finalmente serializados como JSON na resposta HTTP pelo controller.

## 5. Comandos e Scripts

-   **`Makefile`**: Centraliza comandos úteis como `make run`, `make test`, `make migration_up`, e `make load_products`.
-   **Scripts de Carga (`cmd/scripts/load/`)**: Existem scripts para carregar dados legados (provavelmente de um banco MySQL, a julgar pela dependência `go-sql-driver/mysql`) para o novo banco PostgreSQL.

## 6. Pontos de Atenção e Próximos Passos

-   **Soft Delete**: As tabelas `products` e `customers` usam a estratégia de "soft delete" com um campo `is_deleted`. As queries de busca devem sempre incluir a condição `is_deleted = false`. A tabela `addresses` não usa soft delete, sendo removida diretamente através do relacionamento `address_customers`.
-   **Paginação**: A paginação é um requisito central para as listagens (`FindAll`).
-   **Novas Entidades**: Para adicionar uma nova entidade, o fluxo é:
    1.  Criar o `struct` e a interface do repositório em `internal/domain`.
    2.  Criar a migração SQL para a nova tabela.
    3.  Implementar o repositório em `internal/infra/repository/postgres`.
    4.  Criar o `usecase` em `internal/usecase`.
    5.  Criar o `controller` em `internal/controller`.
    6.  Registrar as novas rotas em `cmd/api/main.go`, injetando as dependências.

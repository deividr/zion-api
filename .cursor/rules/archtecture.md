# Cursor Rule: Análise Profunda da Estrutura do Projeto Zion API

## 1. Árvore de Diretórios e Arquivos

```
zion-api/
├── .github/
│   └── workflows/
│       └── fly-deploy.yml
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── scripts/
│       └── load/
│           ├── address/
│           │   └── main.go
│           ├── customers/
│           │   └── main.go
│           └── products/
│               └── main.go
├── data/
│   └── ... (diretórios e arquivos do PostgreSQL)
├── internal/
│   ├── controller/
│   │   ├── address.go
│   │   ├── customer.go
│   │   └── product.go
│   ├── domain/
│   │   ├── address.go
│   │   ├── customer.go
│   │   ├── pagination.go
│   │   ├── product.go
│   │   └── category_product.go
│   ├── infra/
│   │   ├── database/
│   │   │   ├── migrations/
│   │   │   │   └── *.sql
│   │   │   └── postgres.go
│   │   ├── logger/
│   │   │   └── logger.go
│   │   └── repository/
│   │       └── postgres/
│   │           ├── pg_address_repository.go
│   │           ├── pg_customer_repository.go
│   │           └── pg_product_repository.go
│   ├── middleware/
│   │   └── auth.go
│   └── usecase/
│       ├── address.go
│       ├── customer.go
│       └── product.go
├── tmp/
│   └── ... (arquivos temporários)
├── .air.toml
├── .dockerignore
├── .gitignore
├── Dockerfile
├── Dockerfile.develop
├── LICENSE
├── README.md
├── docker-compose.yaml
├── fly.toml
├── go.mod
├── go.sum
└── makefile
```

---

## 2. Explicação das Seções

### Raiz do Projeto

- **README.md**: Documentação principal, visão geral, instruções de uso e arquitetura.
- **Dockerfile / Dockerfile.develop / docker-compose.yaml / fly.toml**: Arquivos de infraestrutura para build, deploy e execução local/produção.
- **go.mod / go.sum**: Gerenciamento de dependências Go.
- **makefile**: Comandos utilitários para build, testes, migrações e importação de dados.

### `.github/`

- **workflows/**: Automação de CI/CD (deploy via Fly.io).

### `cmd/`

- **api/**: Ponto de entrada da API principal (`main.go`).
- **scripts/load/**: Scripts para importação de dados legados (MySQL → PostgreSQL), organizados por domínio (address, customers, products).

### `data/`

- Estrutura de dados do PostgreSQL (usada em dev/local).

### `internal/`

- **controller/**: Camada de controladores HTTP (interface com o framework Gin).
- **domain/**: Entidades de negócio, interfaces de repositório e contratos de domínio.
- **infra/**:
  - **database/**: Conexão e migrações do banco de dados.
  - **logger/**: Implementação de logging estruturado.
  - **repository/postgres/**: Implementações dos repositórios para PostgreSQL.
- **middleware/**: Middlewares HTTP, como autenticação JWT.
- **usecase/**: Casos de uso (aplicação da lógica de negócio).

### `tmp/`

- Arquivos temporários de build e execução.

---

## 3. Principais Pontos de Acoplamento

- **Controllers → Usecases**: Os controladores dependem fortemente dos casos de uso para orquestrar a lógica de negócio.
- **Usecases → Domain Repositories**: Os casos de uso dependem de interfaces de repositório do domínio, promovendo baixo acoplamento com a infraestrutura.
- **Repositórios (infra/repository/postgres) → Banco de Dados**: Implementações concretas dos repositórios acopladas ao PostgreSQL e à biblioteca Squirrel para queries SQL.
- **Middleware de Autenticação → Clerk/JWT**: O middleware depende da chave pública do Clerk para validar tokens JWT.
- **Configuração via Variáveis de Ambiente**: Conexão com banco, autenticação e outros parâmetros críticos dependem de variáveis de ambiente.

---

## 4. Acompanhamentos Eferentes e Aferentes

### Eferentes (para fora do projeto)

- **Frameworks e Bibliotecas**:
  - `gin-gonic/gin` (HTTP)
  - `pgx` (PostgreSQL driver)
  - `squirrel` (query builder)
  - `zerolog` (logging)
  - `golang-jwt/jwt` (autenticação)
  - `joho/godotenv` (env)
  - `clerk` (autenticação externa)
- **Banco de Dados**: PostgreSQL (produção e scripts de importação).
- **MySQL**: Apenas para scripts de importação de dados legados.

### Aferentes (para dentro do projeto)

- **Entradas HTTP**: Todas as requisições passam pelos controllers, que validam, tratam e encaminham para os usecases.
- **Scripts de Importação**: Executam queries em bancos legados e inserem dados via repositórios internos.

---

## 5. Resumo da Estrutura, Propósito, Restrições, Riscos e Segurança

### O que é o projeto?

- **Zion API** é uma API RESTful moderna, escrita em Go, que segue Clean Architecture, focada em gestão de produtos, clientes e endereços, com autenticação JWT via Clerk e persistência em PostgreSQL.

### Restrições

- **Banco de Dados**: Estritamente PostgreSQL (produção), com scripts para migração de dados de MySQL.
- **Autenticação**: Obrigatória via JWT/Clerk para rotas protegidas.
- **Arquitetura**: Segue Clean Architecture, separando domínio, casos de uso, controladores e infraestrutura.
- **Configuração**: Dependente de variáveis de ambiente para funcionamento correto.

### Principais Riscos

- **Configuração de CORS**: Atualmente permite qualquer origem (`AllowOrigins: *`), o que pode ser um risco em produção.
- **Exposição de Variáveis Sensíveis**: Uso de variáveis de ambiente para chaves e conexões, risco se não forem bem gerenciadas.
- **Dependência de Clerk**: Se o serviço de autenticação externo estiver indisponível, a API não autentica usuários.
- **Scripts de Migração**: Scripts de importação podem inserir dados inconsistentes se não houver validação adequada.

### Problemas de Segurança

- **JWT**: A validação depende da chave pública correta do Clerk. Se a variável de ambiente estiver errada ou exposta, pode comprometer a segurança.
- **CORS**: Permitir todas as origens pode expor a API a ataques de CSRF e XSS.
- **Soft Delete**: O uso de `is_deleted` pode permitir acesso a dados "deletados" se não houver filtros corretos em todas as queries.
- **Injeção de SQL**: O uso de Squirrel minimiza, mas não elimina totalmente o risco se usado incorretamente.
- **Logs**: Certifique-se de que informações sensíveis não sejam logadas inadvertidamente.

---

## 6. Resumo Visual de Acoplamentos

```mermaid
graph TD
  subgraph Infraestrutura
    DB[(PostgreSQL)]
    Logger[(Zerolog)]
    Clerk[(Clerk JWT)]
  end

  subgraph Domínio
    Domain[Domain Entities/Interfaces]
  end

  subgraph Aplicação
    Usecase[Usecases]
    Controller[Controllers]
    Middleware[Middleware]
  end

  subgraph Entrada
    API[HTTP API (Gin)]
    Scripts[Import Scripts]
  end

  API --> Middleware
  Middleware --> Controller
  Controller --> Usecase
  Usecase --> Domain
  Usecase -->|Repository| InfraRepo[Infra Repositories]
  InfraRepo --> DB
  Controller --> Logger
  Middleware --> Clerk
  Scripts --> InfraRepo
```

---

## 7. Recomendações

- **CORS**: Restringir origens em produção.
- **Variáveis Sensíveis**: Usar secrets managers e nunca versionar `.env`.
- **Validação de Dados**: Garantir validação rigorosa nos scripts de importação.
- **Auditoria de Logs**: Revisar para evitar vazamento de dados sensíveis.
- **Testes de Segurança**: Realizar pentests e revisões periódicas.

---

**Este documento serve como referência para onboarding, revisão de arquitetura e governança de segurança do projeto Zion API.**

---

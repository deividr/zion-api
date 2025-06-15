# Cursor Rule: Guia de Contexto Completo - Zion API

## 🎯 OBJETIVO DESTA RULE

Esta rule serve como contexto permanente para o Cursor entender completamente o projeto Zion API e gerar código preciso sem necessidade de contextualização adicional.

---

## 📁 ESTRUTURA DO PROJETO

```
zion-api/
├── cmd/
│   ├── api/main.go                    # Entry point da API
│   └── scripts/load/                  # Scripts de migração de dados
├── internal/
│   ├── controller/                    # HTTP handlers (Gin)
│   ├── usecase/                       # Lógica de aplicação
│   ├── domain/                        # Entidades e interfaces
│   ├── infra/
│   │   ├── database/                  # Conexão e migrações
│   │   ├── repository/postgres/       # Implementações de repositório
│   │   └── logger/                    # Logging estruturado
│   └── middleware/                    # Auth, CORS, etc.
├── data/                              # Dados PostgreSQL local
└── [docker/config files]
```

---

## 🏗️ ARQUITETURA E PADRÕES

### Clean Architecture (Obrigatório)

```
HTTP Request → Middleware → Controller → Usecase → Repository → Database
                    ↓           ↓          ↓           ↓
                 Auth/CORS   Validation  Business   Data Layer
```

### Regras de Dependência

- **Controllers** dependem de **Usecases**
- **Usecases** dependem de **Domain interfaces** (não implementações)
- **Domain** é independente de qualquer framework
- **Infra** implementa interfaces do Domain

### Convenções de Nomenclatura

```go
// Domain
type Customer struct {}
type CustomerRepository interface {}

// Usecase
type CustomerUsecase struct {}
func (uc *CustomerUsecase) CreateCustomer() {}

// Controller
type CustomerController struct {}
func (c *CustomerController) CreateCustomer() {}

// Repository
type PgCustomerRepository struct {}
func (r *PgCustomerRepository) Create() {}
```

---

## 🛠️ STACK TECNOLÓGICA

### Core Dependencies

```go
// HTTP Framework
"github.com/gin-gonic/gin"

// Database
"github.com/jackc/pgx/v5"
"github.com/Masterminds/squirrel"

// Auth
"github.com/golang-jwt/jwt/v5"

// Logging
"github.com/rs/zerolog"

// Config
"github.com/joho/godotenv"
```

### Padrões de Código Obrigatórios

#### 1. Estrutura de Controller

```go
type CustomerController struct {
    usecase domain.CustomerUsecase
    logger  zerolog.Logger
}

func (c *CustomerController) CreateCustomer(ctx *gin.Context) {
    var req CreateCustomerRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(400, gin.H{"error": "invalid request"})
        return
    }

    customer, err := c.usecase.CreateCustomer(ctx, req.ToEntity())
    if err != nil {
        c.logger.Error().Err(err).Msg("failed to create customer")
        ctx.JSON(500, gin.H{"error": "internal server error"})
        return
    }

    ctx.JSON(201, NewCustomerResponse(customer))
}
```

#### 2. Estrutura de Usecase

```go
type CustomerUsecase struct {
    repo   domain.CustomerRepository
    logger zerolog.Logger
}

func (uc *CustomerUsecase) CreateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
    // Validações de negócio
    if err := customer.Validate(); err != nil {
        return nil, fmt.Errorf("invalid customer: %w", err)
    }

    // Lógica de aplicação
    existingCustomer, err := uc.repo.FindByEmail(ctx, customer.Email)
    if err != nil && !errors.Is(err, domain.ErrNotFound) {
        return nil, fmt.Errorf("failed to check existing customer: %w", err)
    }

    if existingCustomer != nil {
        return nil, domain.ErrCustomerAlreadyExists
    }

    return uc.repo.Create(ctx, customer)
}
```

#### 3. Estrutura de Repository

```go
type PgCustomerRepository struct {
    db *pgxpool.Pool
}

func (r *PgCustomerRepository) Create(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
    query, args, err := squirrel.
        Insert("customers").
        Columns("name", "email", "created_at").
        Values(customer.Name, customer.Email, time.Now()).
        Suffix("RETURNING id, created_at").
        PlaceholderFormat(squirrel.Dollar).
        ToSql()

    if err != nil {
        return nil, fmt.Errorf("failed to build query: %w", err)
    }

    row := r.db.QueryRow(ctx, query, args...)
    err = row.Scan(&customer.ID, &customer.CreatedAt)
    if err != nil {
        return nil, fmt.Errorf("failed to create customer: %w", err)
    }

    return customer, nil
}
```

#### 4. Estrutura de Domain

```go
// Entidade
type Customer struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    IsDeleted bool      `json:"-"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (c *Customer) Validate() error {
    if c.Name == "" {
        return ErrInvalidName
    }
    if !isValidEmail(c.Email) {
        return ErrInvalidEmail
    }
    return nil
}

// Interface de Repository
type CustomerRepository interface {
    Create(ctx context.Context, customer *Customer) (*Customer, error)
    FindByID(ctx context.Context, id int) (*Customer, error)
    FindByEmail(ctx context.Context, email string) (*Customer, error)
    Update(ctx context.Context, customer *Customer) (*Customer, error)
    Delete(ctx context.Context, id int) error
    List(ctx context.Context, pagination Pagination) ([]*Customer, error)
}

// Errors
var (
    ErrCustomerNotFound      = errors.New("customer not found")
    ErrCustomerAlreadyExists = errors.New("customer already exists")
    ErrInvalidName          = errors.New("invalid customer name")
    ErrInvalidEmail         = errors.New("invalid email format")
)
```

---

## 🔒 AUTENTICAÇÃO E MIDDLEWARE

### JWT Authentication (Clerk)

```go
// Sempre usar este padrão para rotas protegidas
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractBearerToken(c.GetHeader("Authorization"))
        if token == "" {
            c.JSON(401, gin.H{"error": "missing token"})
            c.Abort()
            return
        }

        claims, err := validateJWT(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        c.Set("user_id", claims.Subject)
        c.Next()
    }
}

// Uso nas rotas
protected := router.Group("/api/v1")
protected.Use(AuthMiddleware())
protected.POST("/customers", customerController.CreateCustomer)
```

### CORS Configuration

```go
// SEMPRE configurar CORS restritivo
config := cors.DefaultConfig()
config.AllowOrigins = strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
```

---

## 🗃️ BANCO DE DADOS

### Connection Pattern

```go
// Sempre usar connection pool
func NewPostgresDB(databaseURL string) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to parse database URL: %w", err)
    }

    config.MaxConns = 25
    config.MinConns = 5
    config.MaxConnLifetime = time.Hour
    config.MaxConnIdleTime = time.Minute * 30

    return pgxpool.NewWithConfig(context.Background(), config)
}
```

### Query Patterns com Squirrel

```go
// SELECT
query, args, err := squirrel.
    Select("id", "name", "email", "created_at").
    From("customers").
    Where(squirrel.Eq{"is_deleted": false}).
    Where(squirrel.Eq{"id": id}).
    PlaceholderFormat(squirrel.Dollar).
    ToSql()

// INSERT
query, args, err := squirrel.
    Insert("customers").
    Columns("name", "email", "created_at").
    Values(customer.Name, customer.Email, time.Now()).
    Suffix("RETURNING id, created_at").
    PlaceholderFormat(squirrel.Dollar).
    ToSql()

// UPDATE
query, args, err := squirrel.
    Update("customers").
    Set("name", customer.Name).
    Set("email", customer.Email).
    Set("updated_at", time.Now()).
    Where(squirrel.Eq{"id": customer.ID}).
    Where(squirrel.Eq{"is_deleted": false}).
    PlaceholderFormat(squirrel.Dollar).
    ToSql()

// SOFT DELETE
query, args, err := squirrel.
    Update("customers").
    Set("is_deleted", true).
    Set("updated_at", time.Now()).
    Where(squirrel.Eq{"id": id}).
    PlaceholderFormat(squirrel.Dollar).
    ToSql()
```

### Migrations

- Sempre criar migration files em `internal/infra/database/migrations/`
- Formato: `YYYYMMDDHHMMSS_description.up.sql` e `.down.sql`
- Sempre incluir `is_deleted BOOLEAN DEFAULT FALSE`
- Sempre incluir `created_at` e `updated_at` timestamps

---

## 📝 LOGGING E ERROR HANDLING

### Logging Pattern

```go
// Sempre usar structured logging
logger := zerolog.New(os.Stdout).With().
    Timestamp().
    Str("service", "zion-api").
    Logger()

// No código
logger.Info().
    Str("customer_id", customer.ID).
    Str("action", "create_customer").
    Msg("customer created successfully")

logger.Error().
    Err(err).
    Str("customer_email", customer.Email).
    Msg("failed to create customer")
```

### Error Handling

```go
// Custom errors no domain
var (
    ErrNotFound     = errors.New("resource not found")
    ErrAlreadyExists = errors.New("resource already exists")
    ErrInvalidInput = errors.New("invalid input")
)

// Error handling nos controllers
func (c *Controller) handleError(ctx *gin.Context, err error) {
    switch {
    case errors.Is(err, domain.ErrNotFound):
        ctx.JSON(404, gin.H{"error": "resource not found"})
    case errors.Is(err, domain.ErrAlreadyExists):
        ctx.JSON(409, gin.H{"error": "resource already exists"})
    case errors.Is(err, domain.ErrInvalidInput):
        ctx.JSON(400, gin.H{"error": "invalid input"})
    default:
        c.logger.Error().Err(err).Msg("internal server error")
        ctx.JSON(500, gin.H{"error": "internal server error"})
    }
}
```

---

## 🔧 CONFIGURAÇÃO

### Environment Variables

```go
type Config struct {
    Port           string `env:"PORT" envDefault:"8080"`
    DatabaseURL    string `env:"DATABASE_URL" envRequired:"true"`
    ClerkPublicKey string `env:"CLERK_PEM_PUBLIC_KEY" envRequired:"true"`
    AllowedOrigins string `env:"ALLOWED_ORIGINS" envDefault:"*"`
    LogLevel       string `env:"LOG_LEVEL" envDefault:"info"`
}

// Sempre validar configurações na inicialização
func LoadConfig() (*Config, error) {
    cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    return cfg, nil
}
```

---

## 🧪 TESTES (Quando Solicitado)

### Test Structure

```go
func TestCustomerUsecase_CreateCustomer(t *testing.T) {
    // Setup
    mockRepo := &MockCustomerRepository{}
    logger := zerolog.Nop()
    usecase := NewCustomerUsecase(mockRepo, logger)

    // Test cases
    tests := []struct {
        name     string
        customer *domain.Customer
        mockFn   func(*MockCustomerRepository)
        wantErr  bool
    }{
        {
            name: "success",
            customer: &domain.Customer{
                Name:  "John Doe",
                Email: "john@example.com",
            },
            mockFn: func(m *MockCustomerRepository) {
                m.On("FindByEmail", mock.Anything, "john@example.com").
                    Return(nil, domain.ErrNotFound)
                m.On("Create", mock.Anything, mock.Anything).
                    Return(&domain.Customer{ID: 1}, nil)
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockFn(mockRepo)

            result, err := usecase.CreateCustomer(context.Background(), tt.customer)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, result)
            }
        })
    }
}
```

---

## 📋 REGRAS DE DESENVOLVIMENTO

### 1. **SEMPRE seguir Clean Architecture**

- Domain não conhece infraestrutura
- Usecases orquestram, não implementam regras de dados
- Controllers apenas adaptam HTTP

### 2. **SEMPRE usar interfaces do domain**

```go
// ✅ Correto
type CustomerUsecase struct {
    repo domain.CustomerRepository // Interface
}

// ❌ Errado
type CustomerUsecase struct {
    repo *PgCustomerRepository // Implementação concreta
}
```

### 3. **SEMPRE implementar soft delete**

```sql
-- Todas as tabelas devem ter
is_deleted BOOLEAN DEFAULT FALSE
```

### 4. **SEMPRE usar context.Context**

```go
// ✅ Todas as funções de repository e usecase
func (r *Repository) Create(ctx context.Context, entity *Entity) error

// ✅ Passar contexto do Gin
func (c *Controller) Handler(ginCtx *gin.Context) {
    result, err := c.usecase.DoSomething(ginCtx.Request.Context(), data)
}
```

### 5. **SEMPRE validar entrada**

```go
// Controllers validam formato HTTP
if err := ctx.ShouldBindJSON(&req); err != nil {
    return
}

// Domain valida regras de negócio
if err := entity.Validate(); err != nil {
    return
}
```

### 6. **SEMPRE usar structured logging**

```go
logger.Info().
    Str("entity_id", id).
    Str("action", "operation_name").
    Msg("operation completed")
```

---

## 🚀 PADRÕES DE RESPONSE

### Success Responses

```go
// Single entity
ctx.JSON(200, gin.H{
    "data": entity,
})

// List with pagination
ctx.JSON(200, gin.H{
    "data": entities,
    "pagination": gin.H{
        "page":       page,
        "limit":      limit,
        "total":      total,
        "totalPages": totalPages,
    },
})

// Creation
ctx.JSON(201, gin.H{
    "data":    entity,
    "message": "created successfully",
})
```

### Error Responses

```go
// Client errors (4xx)
ctx.JSON(400, gin.H{"error": "invalid input"})
ctx.JSON(401, gin.H{"error": "unauthorized"})
ctx.JSON(404, gin.H{"error": "not found"})
ctx.JSON(409, gin.H{"error": "already exists"})

// Server errors (5xx)
ctx.JSON(500, gin.H{"error": "internal server error"})
```

---

**IMPORTANTE**: Esta rule deve ser aplicada consistentemente em todo o código. Quando gerar código, SEMPRE siga estes padrões sem exceção. O projeto segue Clean Architecture rigorosamente.

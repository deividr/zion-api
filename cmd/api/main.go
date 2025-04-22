package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	controllers "github.com/deividr/zion-api/internal/controller"
	"github.com/deividr/zion-api/internal/infra/database"
	"github.com/deividr/zion-api/internal/infra/repository/postgres"
	"github.com/deividr/zion-api/internal/middleware"
	"github.com/deividr/zion-api/internal/usecase"
)

func main() {
	godotenv.Load()

	// Setup database connection
	dbPool := database.GetConnection()
	defer dbPool.Close()

	// Setup router
	r := gin.Default()

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	r.Use(cors.New(corsConfig))

	// Grupo de rotas protegidas
	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware(os.Getenv("CLERK_PEM_PUBLIC_KEY")))

	productRoutes(protected, dbPool)
	customerRoutes(protected, dbPool)
	addressRoutes(protected, dbPool)

	r.Run(":8000")
}

func productRoutes(router *gin.RouterGroup, pool *pgxpool.Pool) {
	// Setup repositories
	productRepo := postgres.NewPgProductRepository(pool)

	// Setup use cases
	productUseCase := usecase.NewProductUseCase(productRepo)

	// Setup controllers
	productController := controllers.NewProductController(productUseCase)

	router.GET("/products", productController.GetAll)
	router.GET("/products/:id", productController.GetById)
	router.PUT("/products/:id", productController.Update)
	router.DELETE("/products/:id", productController.Delete)
	router.POST("/products", productController.Create)
}

func customerRoutes(router *gin.RouterGroup, pool *pgxpool.Pool) {
	// Setup repositories
	customerRepo := postgres.NewPgCustomerRepository(pool)
	addressRepo := postgres.NewPgAddressRepository(pool)

	// Setup use cases
	customerUseCase := usecase.NewCustomerUseCase(customerRepo)
	addressUseCase := usecase.NewAddressUseCase(addressRepo)

	// Setup controllers
	customerController := controllers.NewCustomerController(customerUseCase, addressUseCase)

	router.GET("/customers", customerController.GetAll)
	router.GET("/customers/:id", customerController.GetById)
	router.PUT("/customers/:id", customerController.Update)
	router.DELETE("/customers/:id", customerController.Delete)
	router.POST("/customers", customerController.Create)
}

func addressRoutes(router *gin.RouterGroup, pool *pgxpool.Pool) {
	// Setup repositories
	addressRepo := postgres.NewPgAddressRepository(pool)
	addressUseCase := usecase.NewAddressUseCase(addressRepo)

	// Setup controllers
	addressController := controllers.NewAddressController(addressUseCase)

	router.GET("/addresses/:id", addressController.GetById)
	router.PUT("/addresses/:id", addressController.Update)
	router.DELETE("/addresses/:id", addressController.Delete)
	router.POST("/addresses", addressController.Create)
}

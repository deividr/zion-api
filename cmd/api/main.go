package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	// Setup repositories
	productRepo := postgres.NewPgProductRepository(dbPool)

	// Setup use cases
	productUseCase := usecase.NewProductUseCase(productRepo)

	// Setup controllers
	productController := controllers.NewProductController(productUseCase)

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
	{
		protected.GET("/products", productController.GetAll)
		protected.GET("/products/:id", productController.GetById)
		protected.PUT("/products/:id", productController.Update)
		protected.DELETE("/products/:id", productController.Delete)
		protected.POST("/products", productController.Create)
	}

	r.Run(":8000")
}

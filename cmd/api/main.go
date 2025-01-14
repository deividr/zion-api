package main

import (
	controllers "github.com/deividr/zion-api/internal/controller"
	"github.com/deividr/zion-api/internal/infra/database"
	"github.com/deividr/zion-api/internal/infra/repository/postgres"
	"github.com/deividr/zion-api/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
		AllowOrigins:     []string{"*"}, // Altere para os domínios permitidos
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	r.Use(cors.New(corsConfig))

	r.GET("/products", productController.GetAll)
	r.GET("/products/:id", productController.GetById)
	r.PUT("/products/:id", productController.Update)
	r.DELETE("/products/:id", productController.Delete)

	r.Run(":8000")
}

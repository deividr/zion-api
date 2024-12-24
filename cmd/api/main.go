package main

import (
	"github.com/deividr/zion-api/internal/controllers"
	"github.com/deividr/zion-api/internal/infra/database"
	"github.com/deividr/zion-api/internal/infra/repository/postgres"
	"github.com/deividr/zion-api/internal/usecase"
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
	r.GET("/products", productController.GetAll)
	r.GET("/products/:id", productController.GetById)

	r.Run(":8000")
}

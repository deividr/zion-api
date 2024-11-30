package main

import (
	"github.com/deividr/zion-api/internal/controllers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("Environments variables does not allow")
	}

	r := gin.Default()
	r.GET("products", controllers.GetAll)
	r.GET("products/:id", controllers.GetById)
	r.Run(":8000")
}

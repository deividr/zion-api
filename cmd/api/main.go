package main

import (
	"github.com/deividr/zion-api/internal/controllers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	r := gin.Default()
	r.GET("products", controllers.GetAll)
	r.GET("products/:id", controllers.GetById)
	r.Run(":8000")
}

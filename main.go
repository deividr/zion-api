package main

import (
	"github.com/deividr/zion-api/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("products", controllers.GetProducts)
	r.Run("localhost:8000")
}

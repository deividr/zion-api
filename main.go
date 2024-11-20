package main

import (
	"github.com/deividr/zion-api/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("products", controllers.GetAll)
	r.GET("products/:id", controllers.GetById)
	r.Run(":8000")
}

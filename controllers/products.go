package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Product struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Value     uint32 `json:"value"`
	UnityType string `json:"unityType"`
}

var products = []Product{
	{Id: "1", Name: "Blue Train", UnityType: "John Coltrane", Value: 599},
	{Id: "2", Name: "Jeru", UnityType: "Gerry Mulligan", Value: 1799},
	{Id: "3", Name: "Sarah Vaughan and Clifford Brown", UnityType: "Sarah Vaughan", Value: 3999},
}

func GetAll(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, products)
}

func GetById(c *gin.Context) {
	id := c.Param("id")

	for _, product := range products {
		if product.Id == id {
			c.IndentedJSON(http.StatusFound, product)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
}

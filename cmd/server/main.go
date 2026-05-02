package main

import (
	"github.com/gin-gonic/gin"
	"innodeploy-pfe/internal/pipeline"	

)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "INNODEPLOY API running 🚀",
		})
	})

	r.POST("/webhook", func(c *gin.Context) {

		type Request struct {
			Repo string `json:"repo"`
		}

		var req Request

		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		go pipeline.RunPipeline(req.Repo)

		c.JSON(200, gin.H{
			"message": "Pipeline started 🚀",
		})

	})
	r.GET("/webhook", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Webhook GET test",
		})
	})
	r.Run(":8080")
}
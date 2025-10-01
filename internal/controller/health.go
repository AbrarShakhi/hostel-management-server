package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *controller) HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "Hello World",
	})
}

func (h *controller) Health(c *gin.Context) {
	c.JSON(http.StatusOK, h.db.Health())
}

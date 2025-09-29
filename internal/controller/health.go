package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Controller) HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "Hello World",
	})
}

func (h *Controller) Health(c *gin.Context) {
	c.JSON(http.StatusOK, h.db.Health())
}

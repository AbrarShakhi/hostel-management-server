package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Controller) HelloWorld(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (h *Controller) Health(c *gin.Context) {
	c.JSON(http.StatusOK, h.db.Health())
}

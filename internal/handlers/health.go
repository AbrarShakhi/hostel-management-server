package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, h.db.Health())
}

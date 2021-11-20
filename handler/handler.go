package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinHandler struct {
	*gin.Context
}
type SensorHandler interface {
	UpdateSensor()
}
func NewHandler(c *gin.Context) SensorHandler {
	return &GinHandler{
		Context: c,
	}
}
func ( c *GinHandler) UpdateSensor() {
	c.HTML(http.StatusOK, "index", gin.H{
		"title": "Index title!",
		"add": func(a int, b int) int {
			return a + b
		},
	})
}
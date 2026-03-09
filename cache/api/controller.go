package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	cache "distributed-cache/cache/internal"
)

type Handler struct {
	cache *cache.Cache
}

func RegisterRoutes(r *gin.Engine, c *cache.Cache) {
	h := &Handler{cache: c}

	r.POST("/cache/", h.SetValue)
	r.GET("/cache/:key", h.GetValue)
	r.DELETE("/cache/:key", h.DeleteValue)
}

func (h *Handler) SetValue(c *gin.Context) {
	var body struct {
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key := h.cache.SetValue(body.Value)
	c.JSON(http.StatusCreated, gin.H{"key": key})
}

func (h *Handler) GetValue(c *gin.Context) {
	key := c.Param("key")

	value, ok := h.cache.Get(key)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found or expired"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": value})
}

func (h *Handler) DeleteValue(c *gin.Context) {
	key := c.Param("key")

	if _, ok := h.cache.Get(key); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found or expired"})
		return
	}

	h.cache.Delete(key)
	c.Status(http.StatusNoContent)
}

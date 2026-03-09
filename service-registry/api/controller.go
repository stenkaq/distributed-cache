package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	serviceRegistry "distributed-cache/service-registry/internal"
)

type Handler struct {
	service serviceRegistry.ServiceRegistryService
}

func RegisterRoutes(r *gin.Engine, service serviceRegistry.ServiceRegistryService) {
	h := &Handler{
		service: service,
	}

	r.GET("/services/", h.GetService)
	r.POST("/services/", h.AddService)

	r.GET("/services/instances/", h.GetServiceInstance)
	r.POST("/services/instances/", h.AddServiceInstance)
}

func (h *Handler) AddService(c *gin.Context) {
	var body struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	svc, err := h.service.RegisterService(c, body.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": svc.ID})
}

func (h *Handler) GetService(c *gin.Context) {
	var body struct {
		Id string `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	svc, err := h.service.GetService(c, body.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, svc)
}

func (h *Handler) AddServiceInstance(c *gin.Context) {
	var body struct {
		ID     string                         `json:"id"`
		Host   string                         `json:"host"`
		Port   int                            `json:"port"`
		Status serviceRegistry.InstanceStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ServiceId is required"})
		return
	}

	svc, err := h.service.RegisterServiceInstance(c, body.ID, body.Host, body.Port, body.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, svc)
}

func (h *Handler) GetServiceInstance(c *gin.Context) {
	var body struct {
		Id string `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	svc, err := h.service.GetServiceInstance(c, body.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, svc)
}

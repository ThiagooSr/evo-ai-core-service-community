package handler

import (
	"net/http"

	"evo-ai-core-service/internal/httpclient/errors"
	"evo-ai-core-service/internal/httpclient/response"
	"evo-ai-core-service/internal/middleware"
	"evo-ai-core-service/pkg/agent_integration/model"
	"evo-ai-core-service/pkg/agent_integration/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AgentIntegrationHandler interface defines the contract for agent integration handlers
type AgentIntegrationHandler interface {
	RegisterRoutesMiddleware(router gin.IRouter)
	Upsert(c *gin.Context)
	GetByProvider(c *gin.Context)
	ListByAgent(c *gin.Context)
	Delete(c *gin.Context)
}

// agentIntegrationHandler implements the AgentIntegrationHandler interface
type agentIntegrationHandler struct {
	service service.AgentIntegrationService
}

// NewAgentIntegrationHandler creates a new agent integration handler
func NewAgentIntegrationHandler(service service.AgentIntegrationService) AgentIntegrationHandler {
	return &agentIntegrationHandler{
		service: service,
	}
}

// RegisterRoutesMiddleware registers the routes for the agent integration handler with middleware
func (h *agentIntegrationHandler) RegisterRoutesMiddleware(router gin.IRouter) {
	// Get global permission middleware
	permissionMiddleware := middleware.GetGlobalPermissionMiddleware()

	// Routes for agent integrations - using :id instead of :agentId to avoid conflict
	integrations := router.Group("/agents/:id/integrations")
	{
		// List all integrations for an agent
		integrations.GET("",
			permissionMiddleware.RequirePermission("ai_agents", "read"),
			h.ListByAgent)
		integrations.GET("/",
			permissionMiddleware.RequirePermission("ai_agents", "read"),
			h.ListByAgent)

		// Get integration by provider
		integrations.GET("/:provider",
			permissionMiddleware.RequirePermission("ai_agents", "read"),
			h.GetByProvider)

		// Upsert integration
		integrations.POST("",
			permissionMiddleware.RequirePermission("ai_agents", "update"),
			h.Upsert)
		integrations.POST("/",
			permissionMiddleware.RequirePermission("ai_agents", "update"),
			h.Upsert)

		// Delete integration
		integrations.DELETE("/:provider",
			permissionMiddleware.RequirePermission("ai_agents", "update"),
			h.Delete)
	}
}

// Upsert handles the upsert agent integration request
func (h *agentIntegrationHandler) Upsert(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	var req model.AgentIntegrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err)
		return
	}

	integration, err := h.service.Upsert(c.Request.Context(), agentID, req)
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	response.SuccessResponse(c, integration, "Agent integration upserted successfully", http.StatusOK)
}

// GetByProvider handles the get agent integration by provider request
func (h *agentIntegrationHandler) GetByProvider(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	provider := c.Param("provider")
	if provider == "" {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	integration, err := h.service.GetByProvider(c.Request.Context(), agentID, provider)
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	response.SuccessResponse(c, integration, "Agent integration retrieved successfully", http.StatusOK)
}

// ListByAgent handles the list agent integrations request
func (h *agentIntegrationHandler) ListByAgent(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	integrations, err := h.service.ListByAgent(c.Request.Context(), agentID)
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	response.SuccessResponse(c, integrations, "Agent integrations retrieved successfully", http.StatusOK)
}

// Delete handles the delete agent integration request
func (h *agentIntegrationHandler) Delete(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	provider := c.Param("provider")
	if provider == "" {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	err = h.service.Delete(c.Request.Context(), agentID, provider)
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	response.SuccessResponse(c, nil, "Agent integration deleted successfully", http.StatusNoContent)
}

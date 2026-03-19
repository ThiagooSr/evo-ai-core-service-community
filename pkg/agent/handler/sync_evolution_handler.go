package handler

import (
	"net/http"

	"evo-ai-core-service/internal/httpclient/errors"
	"evo-ai-core-service/internal/httpclient/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SyncEvolutionBot handles the sync evolution bot request
func (h *agentHandler) SyncEvolutionBot(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	syncedAgent, err := h.agentService.SyncEvolutionBot(c.Request.Context(), id)
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	response.SuccessResponse(c, syncedAgent.ToResponse(h.aiProcessorURL), "Agent synced with Evolution successfully", http.StatusOK)
}

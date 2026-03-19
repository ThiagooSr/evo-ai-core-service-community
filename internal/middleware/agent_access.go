package middleware

import (
	"evo-ai-core-service/internal/utils/contextutils"
	"net/http"

	agentService "evo-ai-core-service/pkg/agent/service"
	folderShareService "evo-ai-core-service/pkg/folder_share/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AgentAccessMiddleware interface {
	GetAgentAccessMiddleware() gin.HandlerFunc
}

type agentAccessMiddleware struct {
	folderShareService folderShareService.FolderShareService
	agentService       agentService.AgentService
}

func NewAgentAccessMiddleware(folderShareService folderShareService.FolderShareService, agentService agentService.AgentService) AgentAccessMiddleware {
	return &agentAccessMiddleware{folderShareService: folderShareService, agentService: agentService}
}

func bodyToMap(c *gin.Context) (map[string]interface{}, error) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func (a *agentAccessMiddleware) GetAgentAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var requiredPermission string
		var agentID uuid.UUID
		var errAgentID error

		switch c.Request.Method {
		case http.MethodPost:
			requiredPermission = "write"
			agentID, errAgentID = uuid.Parse(c.Param("id"))
		case http.MethodDelete, http.MethodPut:
			requiredPermission = "write"
			agentID, errAgentID = uuid.Parse(c.Param("id"))
		default:
			agentID, errAgentID = uuid.Parse(c.Param("id"))
			requiredPermission = "read"
		}

		_, err := contextutils.GetAccountID(ctx)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
			return
		}

		if errAgentID != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid agent ID"})
			return
		}

		email := ctx.Value("email").(string)

		if email == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
			return
		}

		agent, err := a.agentService.GetByIDAndAccountID(ctx, agentID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid agent ID"})
			return
		}

		if agent.FolderID != nil {
			hasAccess, err := a.folderShareService.CheckFolderAccess(c.Request.Context(), *agent.FolderID, email, requiredPermission)
			if err == nil && hasAccess {
				c.Next()
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Next()
	}
}

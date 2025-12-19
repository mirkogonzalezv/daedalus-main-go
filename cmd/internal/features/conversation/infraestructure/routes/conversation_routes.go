package routes

import (
	authService "daedalus-engine-go/cmd/internal/features/auth/infraestructure/services"
	conversationCtrl "daedalus-engine-go/cmd/internal/features/conversation/infraestructure/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterConversationRoutes(r *gin.RouterGroup, ctrl *conversationCtrl.ConversationController, authService *authService.AuthService) {
	conversations := r.Group("/conversations")

	// Todas las rutas requieren autenticación
	conversations.POST("", ctrl.CreateConversation)
	conversations.GET("", ctrl.ObtenerConversaciones)
	conversations.GET("/:id", ctrl.ObtenerConversacionPorId)
	conversations.PUT("/:id", ctrl.ActualizarConversacion)
	conversations.DELETE("/:id", ctrl.EliminarConversacion)
}

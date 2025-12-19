package routes

import (
	authService "daedalus-engine-go/cmd/internal/features/auth/infraestructure/services"
	controllers "daedalus-engine-go/cmd/internal/features/users/infraestructure/handlers"
	authMiddlewares "daedalus-engine-go/cmd/internal/pkg/middlewares/auth"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, ctr *controllers.UserController, authService *authService.AuthService) {

	user := rg.Group("/users")

	user.POST("", authMiddlewares.AuthMiddleware(authService), authMiddlewares.RequiredGlobalAdmin(), ctr.CrearUsuario)

	user.GET("", authMiddlewares.AuthMiddleware(authService), authMiddlewares.RequiredGlobalAdmin(), ctr.ObtenerUsuario)

	user.GET("/all", authMiddlewares.AuthMiddleware(authService), authMiddlewares.RequiredGlobalAdmin(), ctr.ObtenerListaUsuarios)

	user.GET(":id", authMiddlewares.AuthMiddleware(authService), authMiddlewares.RequiredGlobalAdmin(), ctr.ObtenerUsuarioPorId)

	user.DELETE(":id", authMiddlewares.AuthMiddleware(authService), authMiddlewares.RequiredGlobalAdmin(), ctr.EliminarUsuario)

}

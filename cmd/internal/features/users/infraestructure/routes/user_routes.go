package routes

import (
	authService "daedalus-engine-go/cmd/internal/features/auth/infraestructure/services"
	controllers "daedalus-engine-go/cmd/internal/features/users/infraestructure/handlers"
	middlewares "daedalus-engine-go/cmd/internal/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, ctr *controllers.UserController, authService *authService.AuthService) {

	user := rg.Group("/users")

	user.POST("", middlewares.AuthMiddleware(authService), middlewares.RequiredGlobalAdmin(), ctr.CrearUsuario)

	user.GET("", middlewares.AuthMiddleware(authService), middlewares.RequiredGlobalAdmin(), ctr.ObtenerUsuario)

	user.GET("/all", middlewares.AuthMiddleware(authService), middlewares.RequiredGlobalAdmin(), ctr.ObtenerListaUsuarios)

	user.GET(":id", middlewares.AuthMiddleware(authService), middlewares.RequiredGlobalAdmin(), ctr.ObtenerUsuarioPorId)

	user.DELETE(":id", middlewares.AuthMiddleware(authService), middlewares.RequiredGlobalAdmin(), ctr.EliminarUsuario)

}

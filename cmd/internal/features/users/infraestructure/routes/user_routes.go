package routes

import (
	controllers "daedalus-engine-go/cmd/internal/features/users/infraestructure/handlers"
	middlewares "daedalus-engine-go/cmd/internal/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, ctr *controllers.UserController) {

	user := rg.Group("/users")

	user.POST("", middlewares.AuthMiddleware(), middlewares.RequiredGlobalAdmin(), ctr.CrearUsuario)

	user.GET("", middlewares.AuthMiddleware(), middlewares.RequiredGlobalAdmin(), ctr.ObtenerUsuario)

	user.GET("/all", middlewares.AuthMiddleware(), middlewares.RequiredGlobalAdmin(), ctr.ObtenerListaUsuarios)

	user.GET(":id", middlewares.AuthMiddleware(), middlewares.RequiredGlobalAdmin(), ctr.ObtenerUsuarioPorId)

	user.DELETE(":id", middlewares.AuthMiddleware(), middlewares.RequiredGlobalAdmin(), ctr.EliminarUsuario)

}

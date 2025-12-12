package routes

import (
	"daedalus-engine-go/cmd/core/infraestructure/controllers"
	"daedalus-engine-go/cmd/core/infraestructure/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, ctr *controllers.UserController) {

	user := rg.Group("/users")

	user.POST("", middlewares.AuthMiddleware(), middlewares.RequiredGlobalAdmin(), ctr.CrearUsuario)

	user.GET(":id", middlewares.AuthMiddleware(), middlewares.RequiredGlobalAdmin(), ctr.ObtenerUsuarioPorId)

	user.GET("", middlewares.AuthMiddleware(), middlewares.RequiredGlobalAdmin(), ctr.ObtenerUsuarioPorEmailYTenant)

	user.DELETE(":id", middlewares.AuthMiddleware(), middlewares.RequiredGlobalAdmin(), ctr.EliminarUsuario)
}

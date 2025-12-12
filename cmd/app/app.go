package app

import (
	"daedalus-engine-go/cmd/common/database"
	"daedalus-engine-go/cmd/common/logger"
	"daedalus-engine-go/cmd/config"
	"daedalus-engine-go/cmd/container"
	"daedalus-engine-go/cmd/core/infraestructure/routes"
	"database/sql"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	db     *sql.DB
	router *gin.Engine
	config *config.Config
	log    *zap.Logger
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize() error {
	// Logger
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	config.LoadEnv(env)
	logger.Init(env)
	a.log = logger.L()
	a.log.Info("Daedalus Engine iniciando...", zap.String("env", env))

	// Config
	cfg, err := config.CargarVariables()
	if err != nil {
		return err
	}
	a.config = cfg
	a.log.Info("Configuración cargada correctamente")

	// Database
	db, err := database.NuevaBaseDeDatos(cfg)
	if err != nil {
		return err
	}
	a.db = db
	a.log.Info("Conexión a base de datos establecida")

	// Configurar Gin según entorno
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Dependencies & Routes
	cont := container.NewContainer(db)
	router := gin.New()

	// Agregar middlewares manualmente para evitar warning
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	apiRouter := routes.NewAPIRouter(cont)
	apiRouter.RegisterRouter(router)
	a.router = router

	return nil
}

func (a *App) Run() error {
	a.log.Info("Daedalus Engine iniciado correctamente... ✅")
	return a.router.Run(":" + strconv.Itoa(a.config.Port))
}

func (a *App) Shutdown() {
	if a.db != nil {
		a.db.Close()
	}
}

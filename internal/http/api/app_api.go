package api

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/Intiqo/app-platform/internal/http/handler"
	"github.com/Intiqo/app-platform/internal/pkg/config"
)

type AppApi struct {
	cfg config.AppConfig

	SettingHandler handler.SettingHandler
}

// NewAppApi initializes all the routes for the application.
//
//	@title						App API
//	@version					1.0
//	@description				App's set of APIs
//	@termsOfService				https://intiqo.com/terms/
//	@contact.name				Intiqo Private Limited
//	@contact.url				https://intiqo.in/
//	@contact.email				hello@intiqo.com
//	@host						localhost:8080
//	@BasePath					/api/v1
//	@schemes					https
//	@securityDefinitions.apiKey	JWT
//	@in							header
//
//	@name						Authorization
func NewAppApi(
	cfg config.AppConfig,

	sh handler.SettingHandler,
) *AppApi {
	return &AppApi{
		cfg: cfg,

		SettingHandler: sh,
	}
}

// SetupRoutes initializes all the routes for the application.
func (t AppApi) SetupRoutes(e *echo.Echo) {
	g := e.Group("/api/v1")

	auth := echojwt.JWT([]byte(t.cfg.AuthSecret))

	settingApi := g.Group("/setting")
	settingApi.Use(auth)
	settingApi.GET("/:id", t.SettingHandler.FindByID)
	settingApi.POST("/filter", t.SettingHandler.Filter)
}

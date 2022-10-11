package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/milobella/oratio/internal/config"
)

func ApplyMiddleware(server *echo.Echo, configuration config.Auth) {
	if len(configuration.AppSecret) > 0 {
		// TODO: use custom claim to retrieve scopes and other user info (https://echo.labstack.com/cookbook/jwt)
		//  https://github.com/milobella/oratio/issues/12
		server.Use(middleware.JWT([]byte(configuration.AppSecret)))
	}
}

package handler

import (
	"ahv-worldwide/config"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes wires all domain routes onto the public and admin groups.
func RegisterRoutes(pub, admin *echo.Group, _ *config.Config) {
	registerLeadRoutes(pub, admin)
	registerSettingsRoutes(pub, admin)
}

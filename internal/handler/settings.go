package handler

import (
	"net/http"

	"ahv-worldwide/internal/service"

	"github.com/labstack/echo/v4"
)

func registerSettingsRoutes(pub, admin *echo.Group) {
	pub.GET("/settings", getSettings)
	admin.GET("/settings", getSettings)
	admin.PUT("/settings", updateSettings)
}

func getSettings(c echo.Context) error {
	settings, err := service.GetSettings()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, settings)
}

func updateSettings(c echo.Context) error {
	var body map[string]string
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	if err := service.UpdateSettings(body); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, ok())
}

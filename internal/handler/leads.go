package handler

import (
	"net/http"

	"ahv-worldwide/internal/model"
	"ahv-worldwide/internal/service"

	"github.com/labstack/echo/v4"
)

func registerLeadRoutes(pub, admin *echo.Group) {
	pub.POST("/leads", createLead)
	admin.GET("/leads", getLeads)
	admin.PUT("/leads/:id/status", updateLeadStatus)
	admin.DELETE("/leads/:id", deleteLead)
}

func createLead(c echo.Context) error {
	var l model.Lead
	if err := c.Bind(&l); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	if l.Source == "" {
		l.Source = "contact_form"
	}
	created, err := service.CreateLead(l)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusCreated, created)
}

func getLeads(c echo.Context) error {
	leads, err := service.ListLeads(c.QueryParam("status"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, leads)
}

func updateLeadStatus(c echo.Context) error {
	var body struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err400(err))
	}
	if err := service.UpdateLeadStatus(c.Param("id"), body.Status); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, ok())
}

func deleteLead(c echo.Context) error {
	if err := service.DeleteLead(c.Param("id")); err != nil {
		return c.JSON(http.StatusInternalServerError, err500(err))
	}
	return c.JSON(http.StatusOK, ok())
}

package main

import (
	"log"
	"net/http"

	"ahv-worldwide/config"
	"ahv-worldwide/internal/db"
	"ahv-worldwide/internal/handler"
	mw "ahv-worldwide/internal/middleware"
	"ahv-worldwide/web"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	if err := db.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	// if err := db.Migrate(); err != nil {
	// 	log.Fatalf("migration failed: %v", err)
	// }

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	pub := e.Group("/api")
	admin := e.Group("/api/admin", mw.BasicAuth(cfg.AdminUsername, cfg.AdminPassword))

	handler.RegisterRoutes(pub, admin, cfg)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	handler.RegisterStaticFS(web.Dist)
	e.GET("/*", handler.ServeStatic)

	log.Printf("🚀 AHV Worldwide backend running on :%s", cfg.Port)
	log.Fatal(e.Start(":" + cfg.Port))
}

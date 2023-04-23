package main

import (
	"github.com/Halil-Ibrahim-Kalan/ogrenci-ders-programi/middlewares"
	"github.com/Halil-Ibrahim-Kalan/ogrenci-ders-programi/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	middlewares.Middleware(e)

	// Routes
	routes.Route(e)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

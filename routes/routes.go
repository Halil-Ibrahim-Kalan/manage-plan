package routes

import (
	"net/http"
	"time"

	"github.com/Halil-Ibrahim-Kalan/ogrenci-ders-programi/handlers"

	"github.com/labstack/echo/v4"
)

func Route(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/static/auth/index.html")
	})

	e.GET("/set-cookie/:uid", func(c echo.Context) error {
		cookie := new(http.Cookie)
		cookie.Name = "uid"
		cookie.Value = c.Param("uid")
		cookie.Expires = time.Now().Add(24 * time.Hour)
		cookie.Path = "/"
		cookie.Domain = "localhost"
		cookie.Secure = false
		cookie.HttpOnly = false
		c.SetCookie(cookie)
		return c.Redirect(http.StatusMovedPermanently, "/static/panel/index.html")
	})

	e.Static("/static", "static")
	User(e)
	Plan(e)
}

func User(e *echo.Echo) {
	userRepo := handlers.NewUser()

	e.GET("/api/users/:id", userRepo.GetUserById)
	e.PUT("/api/users", userRepo.UpdateUser)
	e.DELETE("/api/users/:id", userRepo.DeleteUser)

	e.POST("/api/auth/signin", userRepo.SignIn)
	e.POST("/api/auth/signup", userRepo.SignUp)
}

func Plan(e *echo.Echo) {
	planRepo := handlers.NewPlan()

	e.GET("/api/plans", planRepo.GetPlansByUser)
	e.GET("/api/plans/:id", planRepo.GetPlanById)
	e.POST("/api/plans", planRepo.CreatePlanByUser)
	e.DELETE("/api/plans/:id", planRepo.DeletePlan)
	e.PUT("/api/plans/:id", planRepo.UpdatePlan)

	e.GET("/api/plans/week/:page", planRepo.GetWeekly)
	e.GET("/api/plans/month/:page", planRepo.GetMonthly)
	e.GET("/api/plans/count", planRepo.GetNumberOfMountAndWeek)
}

package handlers

import (
	"errors"
	"net/http"
	"time"

	"strconv"

	"github.com/Halil-Ibrahim-Kalan/ogrenci-ders-programi/databases"
	"github.com/Halil-Ibrahim-Kalan/ogrenci-ders-programi/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PlanRepo struct {
	Db *gorm.DB
}
type Count struct {
	Month int
	Week  int
}

func NewPlan() *PlanRepo {
	db := databases.InitDb()
	db.AutoMigrate(&models.Plan{})
	return &PlanRepo{Db: db}
}

func (repository *PlanRepo) checkDateTime(c echo.Context, plan *models.Plan, id int, uid string) (bool, error) {
	format := "15:04"
	t1, err := time.Parse(format, plan.Start)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
		)
		return false, err
	}
	t2, err := time.Parse(format, plan.End)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
		)
		return false, err
	}
	if t1.After(t2) || t1.Equal(t2) {
		c.JSON(
			http.StatusBadRequest,
			ErrorResponse{Code: http.StatusBadRequest, Message: "Do not make plans at this date and time."},
		)
		return false, nil
	}
	times, _ := models.GetDateTime(repository.Db, plan.Date, id, uid)
	for _, _time := range times {
		start, _ := time.Parse(format, _time.Start)
		end, _ := time.Parse(format, _time.End)

		if ((t1.After(start) || t1.Equal(start)) && t1.Before(end)) ||
			(t2.After(start) && (t2.Before(end) || t2.Equal(end)) ||
				(t1.Before(start) && t2.After(end))) {
			c.JSON(
				http.StatusBadRequest,
				ErrorResponse{Code: http.StatusBadRequest, Message: "Do not make plans at this date and time."},
			)
			return false, nil
		}
	}
	return true, nil
}

func (repository *PlanRepo) GetPlanById(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var plan models.Plan
	err := models.GetPlan(repository.Db, &plan, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound,
				ErrorResponse{Code: http.StatusNotFound, Message: ""},
			)
			return nil
		}
		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
		)
		return err
	}
	return c.JSON(http.StatusOK, plan)
}

func (repository *PlanRepo) UpdatePlan(c echo.Context) error {
	var plan models.Plan
	id, _ := strconv.Atoi(c.Param("id"))
	uid := getCookie(c)
	err := models.GetPlan(repository.Db, &plan, id)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
		)
		return err
	}
	c.Bind(&plan)
	plan.UID = uid
	if check, _ := repository.checkDateTime(c, &plan, id, uid); !check {
		return nil
	}
	err = models.UpdatePlan(repository.Db, &plan)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
		)
		return err
	}
	return c.JSON(http.StatusOK, plan)
}

func (repository *PlanRepo) DeletePlan(c echo.Context) error {
	var plan models.Plan
	id, _ := strconv.Atoi(c.Param("id"))
	err := models.DeletePlan(repository.Db, &plan, id)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
		)
		return err
	}
	return c.JSON(http.StatusOK, ErrorResponse{Code: http.StatusOK, Message: "Plan deleted successfully"})
}

func (repository *PlanRepo) GetWeekly(c echo.Context) error {
	var years []models.Year
	var data models.Plan
	c.Bind(&data)
	page, _ := strconv.Atoi(c.Param("page"))
	uid := getCookie(c)
	models.GetYears(repository.Db, &years, uid)
	count := 0
	for _, year := range years {
		var weeks []models.Week
		models.GetWeeksByYear(repository.Db, &weeks, year.Year, uid)
		for _, week := range weeks {
			count++
			if count == page {
				var plans []models.Plan
				models.GetWeekly(repository.Db, &plans, week.Week, year.Year, uid)
				c.JSON(http.StatusOK, plans)
				return nil
			}
		}
	}
	c.JSON(
		http.StatusInternalServerError,
		ErrorResponse{Code: http.StatusInternalServerError, Message: ""},
	)
	return nil
}

func (repository *PlanRepo) GetMonthly(c echo.Context) error {
	var years []models.Year
	var data models.Plan
	c.Bind(&data)
	page, _ := strconv.Atoi(c.Param("page"))
	uid := getCookie(c)
	models.GetYears(repository.Db, &years, uid)
	count := 0
	for _, year := range years {
		var months []models.Month
		models.GetMonthsByYear(repository.Db, &months, year.Year, uid)
		for _, month := range months {
			count++
			if count == page {
				var plans []models.Plan
				models.GetMonthly(repository.Db, &plans, month.Month, year.Year, uid)
				c.JSON(http.StatusOK, plans)
				return nil
			}
		}
	}
	c.JSON(
		http.StatusInternalServerError,
		ErrorResponse{Code: http.StatusInternalServerError, Message: ""},
	)
	return nil
}

func (repository *PlanRepo) GetNumberOfMountAndWeek(c echo.Context) error {
	var years []models.Year
	uid := getCookie(c)
	models.GetYears(repository.Db, &years, uid)
	monthCount := 0
	weekCount := 0
	for _, year := range years {
		var months []models.Month
		models.GetMonthsByYear(repository.Db, &months, year.Year, uid)
		monthCount += len(months)

		var weeks []models.Week
		models.GetWeeksByYear(repository.Db, &weeks, year.Year, uid)
		weekCount += len(weeks)
	}
	return c.JSON(http.StatusOK, Count{Month: monthCount, Week: weekCount})
}

func getCookie(c echo.Context) string {
	cookie, err := c.Cookie("uid")
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/static/auth/index.html")
		return "-1"
	}
	return cookie.Value
}
func (repository *PlanRepo) GetPlansByUser(c echo.Context) error {
	var plans []models.Plan
	uid := getCookie(c)
	err := models.GetPlansByUser(repository.Db, &plans, uid)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
		)
		return err
	}
	return c.JSON(http.StatusOK, plans)
}

func (repository *PlanRepo) CreatePlanByUser(c echo.Context) error {
	uid := getCookie(c)
	if uid == "-1" {
		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: "User id not found"},
		)
		return nil
	}
	var plan models.Plan
	c.Bind(&plan)
	plan.UID = uid
	if check, err := repository.checkDateTime(c, &plan, -1, uid); !check {
		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
		)
		return nil
	}
	createErr := models.CreatePlan(repository.Db, &plan)
	if createErr != nil {
		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: createErr.Error()},
		)
		return createErr
	}
	return c.JSON(http.StatusOK, plan)
}

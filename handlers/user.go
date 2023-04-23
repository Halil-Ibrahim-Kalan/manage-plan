package handlers

import (
	"net/http"

	"errors"
	"strconv"

	"github.com/Halil-Ibrahim-Kalan/ogrenci-ders-programi/databases"
	"github.com/Halil-Ibrahim-Kalan/ogrenci-ders-programi/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserRepo struct {
	Db *gorm.DB
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewUser() *UserRepo {
	db := databases.InitDb()
	db.AutoMigrate(&models.User{})
	return &UserRepo{Db: db}
}

func (repository *UserRepo) GetUserById(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.User
	err := models.GetUserById(repository.Db, &user, id)
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
	return c.JSON(http.StatusOK, user)
}

func (repository *UserRepo) UpdateUser(c echo.Context) error {
	var user models.User
	id, _ := strconv.Atoi(getCookie(c))
	err := models.GetUserById(repository.Db, &user, id)
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
	c.Bind(&user)
	err = models.UpdateUser(repository.Db, &user)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
		)
		return err
	}
	return c.JSON(http.StatusOK, user)
}

func (repository *UserRepo) DeleteUser(c echo.Context) error {
	var user models.User
	id, _ := strconv.Atoi(c.Param("id"))
	err := models.DeleteUser(repository.Db, &user, id)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
		)
		return err
	}
	return c.JSON(http.StatusOK, ErrorResponse{Code: http.StatusOK, Message: "User deleted successfully"})
}

func (repository *UserRepo) SignIn(c echo.Context) error {
	username := c.FormValue("Username")
	password := c.FormValue("Password")
	var user models.User
	err := models.CheckUserPass(repository.Db, &user, username, password)
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return echo.ErrUnauthorized
	}
	return c.Redirect(http.StatusFound, "/set-cookie/"+strconv.Itoa(int(user.ID)))
}

func (repository *UserRepo) SignUp(c echo.Context) error {
	username := c.FormValue("Username")
	check, err := models.CheckUserByUsername(repository.Db, username)
	if err == nil || !check {
		c.Redirect(http.StatusFound, "/")
		return echo.ErrUnauthorized
	}
	var user models.User
	c.Bind(&user)
	createErr := models.CreateUser(repository.Db, &user)
	if createErr != nil {
		c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: createErr.Error()},
		)
		return createErr
	}
	return c.Redirect(http.StatusFound, "/set-cookie/"+strconv.Itoa(int(user.ID)))
}

package handler

import (
	"backend-github-trending/log"
	"backend-github-trending/model"
	req "backend-github-trending/model/req"
	"backend-github-trending/repository"
	"backend-github-trending/security"
	validator "github.com/go-playground/validator/v10"
	uuid "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserHandler struct {
	UserRepo repository.UserRepo
}

func (u * UserHandler) HandleSignIn(c echo.Context) error {
	req := req.RequestSignIn{}
	if err := c.Bind(&req); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message: err.Error(),
			Data: nil,
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message: err.Error(),
			Data: nil,
		})
	}

	user, err := u.UserRepo.CheckLogin(c.Request().Context(), req)
	if err != nil{
		log.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message: err.Error(),
			Data: nil,
		})
	}

	// CheckPass
	isTheSame := security.ComparePasswords(user.Password, []byte(req.Password))
	if !isTheSame{
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message: "Login Fail",
			Data: nil,
		})
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message: "Success",
		Data: user,
	})
}

func (u * UserHandler) HandleSignUp(c echo.Context) error {
	req := req.RequestSignUp{}
	if err := c.Bind(&req); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message: err.Error(),
			Data: nil,
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message: err.Error(),
			Data: nil,
		})
	}

	hash := security.HashAndSalt([]byte(req.Password))
	role := model.MEMBER.String()

	userId, err := uuid.NewUUID()
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
			Message: err.Error(),
			Data: nil,
		})
	}

	user := model.User{
		UserId: userId.String(),
		FullName: req.FullName,
		Email: req.Email,
		Password: hash,
		Role: role,
		Token: "",
	}
	user, err = u.UserRepo.SaveUser(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusConflict, model.Response{
			StatusCode: http.StatusConflict,
			Message: err.Error(),
			Data: nil,
		})
	}

	user.Password = ""

	return c.JSON(http.StatusOK,  model.Response{
		StatusCode: http.StatusOK,
		Message: "Success",
		Data: user,
	})
}
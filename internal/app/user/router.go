package user

import (
	"errors"
	"net/http"
	"note-service/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
	"note-service/internal/app"
	userpkg "note-service/internal/pkg/user"
)

type userService interface {
	SignUp(name, password string) (userpkg.User, error)
	Login(name, password string) (userpkg.User, error)
}

type Router struct {
	service userService
}

func NewRouter(service userService) *Router {
	return &Router{service: service}
}

func (r *Router) SetUpRouter(engine *gin.Engine) {
	engine.POST("/user", r.signUp)
	engine.POST("/user/login", r.login)
}

func (r *Router) signUp(c *gin.Context) {
	var request SignUpRequest
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, app.ErrorModel{Error: err.Error()})
		return
	}
	err := request.Validate()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}
	u, err := r.service.SignUp(request.Username, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, userpkg.ErrEmptyPassword):
			c.IndentedJSON(http.StatusBadRequest, app.ErrorModel{Error: err.Error()})
		case errors.Is(err, userpkg.ErrUsedUsername):
			c.IndentedJSON(http.StatusConflict, app.ErrorModel{Error: err.Error()})
		default:
			c.IndentedJSON(http.StatusInternalServerError, app.UnknownError)
		}
		return
	}
	c.IndentedJSON(http.StatusCreated, userToUserResponse(u))
}

func (r *Router) login(c *gin.Context) {
	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, app.ErrorModel{Error: err.Error()})
		return
	}
	err := request.Validate()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}
	u, err := r.service.Login(request.Username, request.Password)
	if err != nil {
		if errors.Is(err, userpkg.ErrUserNotFound) {
			c.IndentedJSON(http.StatusNotFound, app.ErrorModel{Error: err.Error()})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, app.UnknownError)
		}
		return
	}

	token, err := jwt.CreateToken(u.ID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, app.ErrorModel{Error: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, app.TokenModel{Token: token})
}

package user

import (
	"errors"
	"go.uber.org/zap"
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
	logger  *zap.Logger
}

func NewRouter(service userService, logger *zap.Logger) *Router {
	return &Router{service: service, logger: logger}
}

func (r *Router) SetUpRouter(engine *gin.Engine) {
	engine.POST("/user", r.signUp)
	engine.POST("/user/login", r.login)
}

func (r *Router) signUp(c *gin.Context) {
	var request SignUpRequest
	if err := c.BindJSON(&request); err != nil {
		r.logger.Error("failed to bind json", zap.Error(err))
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
		if errors.Is(err, userpkg.ErrUsedUsername) {
			c.IndentedJSON(http.StatusConflict, app.ErrorModel{Error: err.Error()})
		} else {
			r.logger.Error("failed to create jwt-token", zap.Error(err))
			c.IndentedJSON(http.StatusInternalServerError, app.UnknownError)
		}
		return
	}
	r.logger.Info("user was created", zap.Any("user", userToUserResponse(u)))
	c.IndentedJSON(http.StatusCreated, userToUserResponse(u))
}

func (r *Router) login(c *gin.Context) {
	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		r.logger.Error("failed to bind json", zap.Error(err))
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
			r.logger.Error("failed to create jwt-token", zap.Error(err))
			c.IndentedJSON(http.StatusInternalServerError, app.UnknownError)
		}
		return
	}

	token, err := jwt.CreateToken(u.ID)
	if err != nil {
		r.logger.Error("failed to create jwt-token", zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, app.ErrorModel{Error: err.Error()})
		return
	}
	r.logger.Info("user was authorized")
	c.IndentedJSON(http.StatusOK, app.TokenModel{Token: token})
}

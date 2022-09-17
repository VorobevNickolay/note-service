package user

import (
	"context"
	"errors"
	"net/http"
	"note-service/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
	"note-service/internal/app"
	userpkg "note-service/internal/pkg/user"
)

type userStore interface {
	CreateUser(ctx context.Context, name, password string) (userpkg.User, error)
	FindUserByID(ctx context.Context, id string) (userpkg.User, error)
	FindUserByNameAndPassword(ctx context.Context, name, password string) (userpkg.User, error)
}

type Router struct {
	store userStore
}

func NewRouter(store userStore) *Router {
	return &Router{store}
}

func (r *Router) SetUpRouter(engine *gin.Engine) {
	engine.GET("/user/:id", r.getUserByID)
	engine.POST("/user", r.signUp)
	engine.POST("/user/login", r.login)
}

func (r *Router) getUserByID(c *gin.Context) {
	id := c.Param("id")
	u, err := r.store.FindUserByID(c, id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, app.ErrorModel{Error: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, userModelFromUser(u))
}

func (r *Router) signUp(c *gin.Context) {
	var newUser userpkg.User
	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, app.ErrorModel{Error: err.Error()})
		return
	}

	u, err := r.store.CreateUser(c, newUser.Username, newUser.Password)
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
	c.IndentedJSON(http.StatusCreated, userModelFromUser(u))
}

func (r *Router) login(c *gin.Context) {
	var u userpkg.User

	if err := c.ShouldBindJSON(&u); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, app.ErrorModel{Error: err.Error()})
		return
	}
	u, err := r.store.FindUserByNameAndPassword(c, u.Username, u.Password)
	if err != nil {
		if errors.Is(err, userpkg.ErrUserNotFound) {
			c.IndentedJSON(http.StatusNotFound, app.ErrorModel{Error: err.Error()})
		} else if errors.Is(err, userpkg.ErrEmptyPassword) {
			c.IndentedJSON(http.StatusBadRequest, app.ErrorModel{Error: err.Error()})
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

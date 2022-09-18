package app

import (
	"errors"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ErrNoAccess = errors.New("you have no access for this action")

type Router struct {
	ginContext *gin.Engine
	subRouters []subRouter
}
type subRouter interface {
	SetUpRouter(engine *gin.Engine)
}

func NewRouter(logger *zap.Logger, subRouters ...subRouter) *Router {
	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	return &Router{r, subRouters}
}

func (r *Router) SetUpRouter() {
	for _, s := range r.subRouters {
		s.SetUpRouter(r.ginContext)
	}
}

func (r *Router) Run() {
	_ = r.ginContext.Run("localhost:8080")
}

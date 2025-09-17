package types

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IKVStore interface {
	Put(c *gin.Context, logger *zap.SugaredLogger)
	Delete(c *gin.Context, logger *zap.SugaredLogger)
	Get(c *gin.Context, logger *zap.SugaredLogger)
}

type IServiceRunner interface {
	Start() error
}

// error declarations, can be extended. For now, should just be element not found in KV
var (
	ErrElementNotFound = errors.New("requested element not found")
)

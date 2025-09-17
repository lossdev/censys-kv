package logic

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lossdev/censys-kv/kv-service/types"
	"go.uber.org/zap"
)

type kvStore struct {
	store map[string]string
	lock  *sync.Mutex
}

func NewKVStore(kv map[string]string, lock *sync.Mutex) types.IKVStore {
	return &kvStore{kv, lock}
}

type serviceRunner struct {
	ginRouter *gin.Engine
	logger    *zap.SugaredLogger
	kvStore   types.IKVStore
}

func NewServiceRunner(ginRouter *gin.Engine, logger *zap.SugaredLogger, kv types.IKVStore) types.IServiceRunner {
	return &serviceRunner{ginRouter, logger, kv}
}

func (sr serviceRunner) Start() error {
	sr.logger.Infoln("Starting kv-service on :8080")

	// route defs
	sr.ginRouter.PUT("/key/:key/:value", func(c *gin.Context) {
		sr.kvStore.Put(c, sr.logger)
	})

	sr.ginRouter.DELETE("/key/:key", func(c *gin.Context) {
		sr.kvStore.Delete(c, sr.logger)
	})

	sr.ginRouter.GET("/key/:key", func(c *gin.Context) {
		sr.kvStore.Get(c, sr.logger)
	})

	sr.ginRouter.Run("0.0.0.0:8080")
	return nil
}

func (kv kvStore) Put(c *gin.Context, logger *zap.SugaredLogger) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	key := c.Params.ByName("key")
	value := c.Params.ByName("value")
	_, ok := kv.store[key]
	kv.store[key] = value
	if ok {
		c.JSON(http.StatusCreated, gin.H{"status": "updated", "key": key, "value": value})
	} else {
		c.JSON(http.StatusCreated, gin.H{"status": "created", "key": key, "value": value})
	}
	logger.Infof("PUT kv[%s] = %s", key, value)
}

func (kv kvStore) Delete(c *gin.Context, logger *zap.SugaredLogger) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	key := c.Params.ByName("key")
	_, ok := kv.store[key]
	if ok {
		delete(kv.store, key)
		c.JSON(http.StatusOK, gin.H{"status": "deleted", "key": key})
		logger.Infof("DELETE kv[%s]", key)
	} else {
		c.JSON(http.StatusOK, gin.H{"message": types.ErrElementNotFound.Error(), "key": key})
	}
}

func (kv kvStore) Get(c *gin.Context, logger *zap.SugaredLogger) {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	key := c.Params.ByName("key")
	value, ok := kv.store[key]
	if ok {
		c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
		logger.Infof("GET kv[%s] = %s", key, value)
	} else {
		c.JSON(http.StatusOK, gin.H{"message": types.ErrElementNotFound.Error(), "key": key})
		logger.Infof("GET kv[%s] not found", key)
	}
}

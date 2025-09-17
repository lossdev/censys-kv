package main

import (
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lossdev/censys-kv/kv-service/internal/logger"
	"github.com/lossdev/censys-kv/kv-service/internal/logic"
)

func main() {
	logger := logger.NewLogger()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	lock := &sync.Mutex{}
	kvMap := make(map[string]string)
	kv := logic.NewKVStore(kvMap, lock)
	svc := logic.NewServiceRunner(r, logger, kv)
	log.Fatalln(svc.Start())
}

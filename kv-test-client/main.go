package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gojek/heimdall/v7/httpclient"
	"go.uber.org/zap"
)

func main() {
	logger := newLogger()
	r := setupRouter()
	logger.Infoln("Running client on :8081, ready to accept /test_deletion and /test_overwrite")
	logger.Fatalln(r.Run("0.0.0.0:8081"))
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/test_deletion", func(c *gin.Context) {
		getTestDeletion(c)
	})

	r.GET("/test_overwrite", func(c *gin.Context) {
		getTestOverwrite(c)
	})

	return r
}

func getTestDeletion(c *gin.Context) {
	client := httpclient.NewClient()
	res, err := client.Put("http://server:8080/key/foo/bar", nil, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed ❌", "message": "unexpected error calling PUT request: " + err.Error()})
	}

	if res.StatusCode != http.StatusCreated {
		c.JSON(http.StatusOK, gin.H{"status": "failed ❌", "message": "expected PUT 201 response, received " + strconv.Itoa(res.StatusCode) + " instead"})
	}

	res2, err := client.Delete("http://server:8080/key/foo", nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed ❌", "message": "unexpected error calling DELETE request: " + err.Error()})
	}
	delResp := &DeleteResponseShape{}
	getJson(res2.Body, delResp)
	if res2.StatusCode != http.StatusOK || delResp.Status != "deleted" {
		c.JSON(http.StatusOK, gin.H{"status": "failed ❌", "message": "key was not deleted"})
	}
	c.JSON(http.StatusOK, gin.H{"status": "passed ✅", "message": "key 'foo' was successfully created and deleted"})
}

func getTestOverwrite(c *gin.Context) {
	client := httpclient.NewClient()
	res, err := client.Put("http://server:8080/key/foo/bar", nil, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed ❌", "message": "unexpected error calling PUT request: " + err.Error()})
	}

	if res.StatusCode != http.StatusCreated {
		c.JSON(http.StatusOK, gin.H{"status": "failed ❌", "message": "expected PUT 201 response, received " + strconv.Itoa(res.StatusCode) + " instead"})
	}

	res2, err := client.Get("http://server:8080/key/foo", nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed ❌", "message": "unexpected error calling PUT request: " + err.Error()})
	}
	getResp := &GetResponseShape{}
	getJson(res2.Body, getResp)
	if res2.StatusCode != http.StatusOK || getResp.Value != "bar" {
		c.JSON(http.StatusOK, gin.H{"status": "failed ❌", "message": "key was not found"})
	}

	res3, err := client.Put("http://server:8080/key/foo/baz", nil, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed ❌", "message": "unexpected error calling PUT request: " + err.Error()})
	}
	if res3.StatusCode != http.StatusCreated {
		c.JSON(http.StatusOK, gin.H{"status": "failed ❌", "message": "expected PUT 201 response, received " + strconv.Itoa(res.StatusCode) + " instead"})
	}

	res4, err := client.Get("http://server:8080/key/foo", nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed ❌", "message": "key was not found: " + err.Error()})
	}
	getResp2 := &GetResponseShape{}
	getJson(res4.Body, getResp2)
	if res4.StatusCode != http.StatusOK || getResp2.Value != "baz" {
		c.JSON(http.StatusOK, gin.H{"status": "failed ❌", "message": "key was not found, recieved " + getResp2.Value})
	}
	c.JSON(http.StatusOK, gin.H{"status": "passed ✅", "message": "key 'foo' was created with value 'bar', updated to 'baz', and correctly overwritten"})
}

func newLogger() *zap.SugaredLogger {
	zapConfig := zap.NewProductionConfig()
	zapConfig.DisableStacktrace = true
	zapConfig.EncoderConfig.CallerKey = ""
	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatalln(err)
	}
	return logger.Sugar()
}

func getJson(r io.ReadCloser, target interface{}) error {
	defer r.Close()
	return json.NewDecoder(r).Decode(target)
}

type DeleteResponseShape struct {
	Status string `json:"status"`
	Key    string `json:"key"`
}

type GetResponseShape struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

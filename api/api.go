package api

import (
	"bytes"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	once   sync.Once
	router *gin.Engine
	root   *gin.RouterGroup
)

func GetRouter() *gin.Engine {
	once.Do(initializeSingletons)
	return router
}

func GetRoot() *gin.RouterGroup {
	once.Do(initializeSingletons)
	return root
}

func initializeSingletons() {
	router, root = createRouterAndGroup("")
}

func createRouterAndGroup(prefix string) (*gin.Engine, *gin.RouterGroup) {
	engine := gin.New()

	engine.RedirectTrailingSlash = true
	engine.RedirectFixedPath = false
	engine.HandleMethodNotAllowed = false
	engine.ForwardedByClientIP = true

	engine.Use(printAllBody)

	group := engine.Group(prefix)

	installCommonMiddleware(group)

	return engine, group
}

func printAllBody(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	println(string(body))

	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	c.Next()
}

func installCommonMiddleware(group *gin.RouterGroup) {
	group.Use(gin.Logger())
	group.Use(gin.Recovery())
}

package main

import "github.com/gin-gonic/gin"

func newEngine() *gin.Engine {

	engine := gin.Default()

	return engine
}

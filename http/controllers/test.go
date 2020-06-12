package controllers

import (
	"github.com/gin-gonic/gin"
	"tcp_client/providers"
)

type Test struct {
	Controller
}

func (test *Test)TestA(c *gin.Context)  {
	providers.Logger.Info("test12312312")
}
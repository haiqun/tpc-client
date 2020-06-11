package middles

import (
	"github.com/gin-gonic/gin"
)

func EmptyMiddles(context *gin.Context) {
	context.Next()
}

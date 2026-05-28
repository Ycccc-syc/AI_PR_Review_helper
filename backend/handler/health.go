package handler

import (
	"qiniudemo/model"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	model.Success(c, gin.H{"status": "ok"})
}

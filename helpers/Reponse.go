package helpers

import (
	"github.com/gin-gonic/gin"
)

type ReponseData struct {
	Status  int
	Message string
	Data    interface{}
}

func RespondJSON(w *gin.Context, status int, message string, payload interface{}) {
	var res ReponseData
	res.Status = status
	res.Message = message
	res.Data = payload

	w.JSON(200, res)
}

package helpers

import (
	"github.com/gin-gonic/gin"
)

type ReponseData struct {
	Status  int
	Message string
	Error   interface{}
	Data    interface{}
}

func RespondJSON(w *gin.Context, status int, message string, errors interface{}, payload interface{}) {
	var res ReponseData
	res.Status = status
	res.Message = message
	res.Error = errors
	res.Data = payload

	w.JSON(200, res)
}

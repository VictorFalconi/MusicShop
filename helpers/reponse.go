package helpers

import (
	"github.com/gin-gonic/gin"
)

// ResponseData : Response JSON data to client
type ResponseData struct {
	Status  int
	Message string
	Error   interface{}
	Data    interface{}
}

func RespondJSON(w *gin.Context, status int, message string, errors interface{}, payload interface{}) {
	var res ResponseData
	res.Status = status
	res.Message = message
	res.Error = errors
	res.Data = payload

	w.JSON(status, res)
}

// LineError : return position & fields when don't create elements from file (excel, csv)
type LineError struct {
	Line    int
	Message interface{}
}

package handler

import (
	"net/http"
	"router-template/delivery"
	"router-template/entities/app"
	"router-template/entities/statuscode"
	"router-template/usecase"

	"github.com/gin-gonic/gin"
)

func GetEmployeListHandler(ctx *gin.Context) {
	ucase := usecase.NewEmployeeUsecase()
	data, er := ucase.GetEmployeeList()
	if er != nil {
		if er == app.ErrNoRecord {
			ctx.String(statuscode.StatusNoRecord, "record not found.")
		} else {
			delivery.PrintError(er.Error())
			ctx.String(http.StatusInternalServerError, "internal service error")
		}
	} else {
		ctx.JSON(http.StatusOK, data)
	}
}

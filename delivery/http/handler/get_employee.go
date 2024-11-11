package handler

import (
	"net/http"
	"router-template/delivery"
	"router-template/entities"
	"router-template/entities/app"
	"router-template/entities/statuscode"
	"router-template/usecase"

	"github.com/gin-gonic/gin"
)

func GetEmployee(ctx *gin.Context) {
	payload := entities.EmployeeFilter{}
	var er error
	if ctx.ContentType() == "application/json" {
		er = ctx.BindJSON(&payload)
	} else {
		er = ctx.Bind(&payload)
	}
	if er != nil {
		ctx.String(http.StatusBadRequest, er.Error())
	}

	ucase := usecase.NewEmployeeUsecase()
	employee, er := ucase.GetEmployee(payload.Id)
	if er != nil {
		if er == app.ErrDuplicateEntry {
			ctx.String(statuscode.StatusDuplicate, "Data karyawan sudah tersedia!")
		} else {
			delivery.PrintError(er.Error())
			ctx.String(http.StatusInternalServerError, "internal service error")
		}
	} else {
		ctx.JSON(http.StatusOK, employee)
	}
}

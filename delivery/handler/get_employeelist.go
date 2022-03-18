package handler

import (
	"clean-arch-employee/delivery/handler/httpio"
	"clean-arch-employee/entities"
	"clean-arch-employee/entities/err"
	"clean-arch-employee/entities/statuscode"
	"clean-arch-employee/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetEmployeListHandler(ctx *gin.Context) {
	httpio := httpio.NewRequestIO(ctx)
	httpio.Recv()

	ucase := usecase.NewEmployeeUsecase()
	data, er := ucase.GetEmployeeList()
	if er != nil {
		if er == err.NoRecord {
			httpio.ResponseString(statuscode.StatusNoRecord, "record not found.")
		} else {
			entities.PrintError(er.Error())
			httpio.ResponseString(http.StatusInternalServerError, "internal service error")
		}
	} else {
		httpio.Response(http.StatusOK, data)
	}
}

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

func GetEmployee(ctx *gin.Context) {
	httpio := httpio.NewRequestIO(ctx)

	payload := entities.EmployeeFilter{}
	httpio.BindJSON(&payload)

	ucase := usecase.NewEmployeeUsecase()
	employee, er := ucase.GetEmployee(int(payload.Id))
	if er != nil {
		if er == err.DuplicateEntry {
			httpio.ResponseString(statuscode.StatusDuplicate, "Data karyawan sudah tersedia!")
		} else {
			entities.PrintError(er.Error())
			httpio.ResponseString(http.StatusInternalServerError, "internal service error")
		}
	} else {
		httpio.Response(http.StatusOK, employee)
	}
}

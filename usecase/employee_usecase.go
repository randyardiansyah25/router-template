package usecase

import (
	"router-template/entities"
	"router-template/entities/app"
	"router-template/repository/employeerepo"
)

type EmployeeUsecase interface {
	GetEmployeeList() ([]entities.Employee, error)
	GetEmployee(id int64) (entities.Employee, error)
}

func NewEmployeeUsecase() EmployeeUsecase {
	return &employeeUsecase{}
}

type employeeUsecase struct{}

func (e *employeeUsecase) GetEmployeeList() (detail []entities.Employee, er error) {
	repo, _ := employeerepo.NewEmployeeRepo()
	detail, er = repo.GetEmployee()

	if er != nil {
		return detail, er
	}

	// ! DO ANOTHER BUSINESS LOGIC HERE

	if len(detail) == 0 {
		return detail, app.ErrNoRecord
	}

	return
}

func (e *employeeUsecase) GetEmployee(id int64) (employee entities.Employee, er error) {
	repo, _ := employeerepo.NewEmployeeRepo()
	employee, er = repo.GetEmployeeById(id)

	// ! DO ANOTHER BUSINESS LOGIC HERE

	return
}

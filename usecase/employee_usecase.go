package usecase

import (
	"clean-arch-employee/entities"
	"clean-arch-employee/entities/err"
	"clean-arch-employee/repository/employeerepo"
)

type EmployeeUsecase interface {
	GetEmployeeList() ([]entities.Employee, error)
	GetEmployee(id int) (entities.Employee, error)
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
		return detail, err.NoRecord
	}

	return
}

func (e *employeeUsecase) GetEmployee(id int) (employee entities.Employee, er error) {
	repo, _ := employeerepo.NewEmployeeRepo()
	employee, er = repo.GetEmployeeById(id)

	// ! DO ANOTHER BUSINESS LOGIC HERE

	return
}

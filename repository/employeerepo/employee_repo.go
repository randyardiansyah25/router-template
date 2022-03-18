package employeerepo

import (
	"clean-arch-employee/entities"
)

type EmployeeRepo interface {
	GetEmployee() ([]entities.Employee, error)
	GetEmployeeById(id int) (entities.Employee, error)
}

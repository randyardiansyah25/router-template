package employeerepo

import (
	"clean-arch-employee/entities"
	"clean-arch-employee/entities/err"
	"database/sql"
	"errors"
	"fmt"
)

func newEmployeeMysqlImpl(conn *sql.DB) EmployeeRepo {
	return &employeeMysqlImpl{conn: conn}
}

type employeeMysqlImpl struct {
	conn *sql.DB
}

func (e *employeeMysqlImpl) GetEmployee() (list []entities.Employee, er error) {
	rows, er := e.conn.Query("Select id, name, address, phone_number")
	if er != nil {
		return list, er
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var item entities.Employee
		if er = rows.Scan(&item.Id, &item.Name, &item.Address, &item.PhoneNumber); er != nil {
			return list, er
		}

		list = append(list, item)
	}

	if len(list) == 0 {
		return list, err.NoRecord
	} else {
		return
	}
}

func (e *employeeMysqlImpl) GetEmployeeById(id int) (employee entities.Employee, er error) {
	row := e.conn.QueryRow(`SELECT 
		id,
		name,
		address,
		phone_number 
		FROM employee WHERE id=?`, id)

	if er = row.Scan(&employee.Id, &employee.Name, &employee.Address, &employee.PhoneNumber); er != nil {
		if er == sql.ErrNoRows {
			return
		} else {
			return employee, errors.New(fmt.Sprint("error while get employee : ", er.Error()))
		}
	}

	return
}

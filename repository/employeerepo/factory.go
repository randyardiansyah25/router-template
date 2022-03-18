package employeerepo

import (
	"clean-arch-employee/repository/databasefactory"
	"clean-arch-employee/repository/databasefactory/drivers"
	"database/sql"
	"errors"
)

func NewEmployeeRepo() (EmployeeRepo, error) {
	conn := databasefactory.AppDb.GetConnection()
	currentDriver := databasefactory.AppDb.GetDriverName()
	if currentDriver == drivers.MYSQL {
		return newEmployeeMysqlImpl(conn.(*sql.DB)), nil
	} else {
		return nil, errors.New("unimplemented database driver")
	}

}

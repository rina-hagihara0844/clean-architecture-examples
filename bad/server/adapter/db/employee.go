package db

import (
	"database/sql"

	"example.com/server/domain"
)

func GetEmployeeByID(db *sql.DB, employeeID string) (domain.Employee, error) {
	var emp domain.Employee
	err := db.QueryRow(`SELECT id, name, hire_date FROM employees WHERE id=$1`, employeeID).Scan(&emp)
	return emp, err
}

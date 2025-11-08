package db

import (
	"database/sql"
	"time"

	"example.com/server/domain"
)

func CountSubmitThisFiscalYear(db *sql.DB, employeeID string) (int, error) {
	yearStart := time.Date(time.Now().Year(), 4, 1, 0, 0, 0, 0, time.UTC)
	var cnt int
	err := db.QueryRow(`SELECT COUNT(*) FROM leave_requests WHERE employee_id=$1 AND created_at>= $2 AND status IN ('PENDING','APPROVED','REJECTED','RETURNED')`,
		employeeID, yearStart).Scan(&cnt)
	return cnt, err
}

func RegisterLeaveRequest(db *sql.DB, submit domain.Submit) (string, error) {
	var id string
	err := db.QueryRow(
		`INSERT INTO leave_requests(employee_id,reason,from_date,to_date,status,created_at)
			 VALUES($1,$2,$3,$4,'PENDING',NOW()) RETURNING id`,
		submit.EmployeeID, submit.Reason, submit.From, submit.To).Scan(&id)
	return id, err
}

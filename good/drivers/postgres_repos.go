package drivers

import (
	"database/sql"
	"time"

	"github.com/ohagi/clean-architecture-examples/good/domain"
)

type PostgresEmployeeRepo struct{ DB *sql.DB }

func (r PostgresEmployeeRepo) FindByID(id string) (domain.Employee, error) {
	var e domain.Employee
	return e, r.DB.QueryRow(`SELECT id, hire_date FROM employees WHERE id=$1`, id).Scan(&e.ID, &e.HireDate)
}

type PostgresLeaveRepo struct{ DB *sql.DB }

func (r PostgresLeaveRepo) CountThisFiscalYear(empID string, start time.Time) (int, error) {
	var c int
	return c, r.DB.QueryRow(
		`SELECT COUNT(*) FROM leave_requests WHERE employee_id=$1 AND created_at >= $2`,
		empID, start).Scan(&c)
}
func (r PostgresLeaveRepo) Create(req *domain.LeaveRequest) error {
	return r.DB.QueryRow(
		`INSERT INTO leave_requests(employee_id,reason,from_date,to_date,status,created_at)
		 VALUES($1,$2,$3,$4,$5,$6) RETURNING id`,
		req.EmployeeID, req.Reason, req.From, req.To, req.Status, req.CreatedAt,
	).Scan(&req.ID)
}

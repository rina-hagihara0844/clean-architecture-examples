package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/ohagi/clean-architecture-examples/good/adapters"
	"github.com/ohagi/clean-architecture-examples/good/drivers"
	"github.com/ohagi/clean-architecture-examples/good/usecase"
)

type sysClock struct{}

func (sysClock) Now() time.Time { return time.Now().UTC() }
func fiscalYearStart(now time.Time) time.Time {
	year := now.Year()
	start := time.Date(year, 4, 1, 0, 0, 0, 0, time.UTC)
	if now.Before(start) {
		start = time.Date(year-1, 4, 1, 0, 0, 0, 0, time.UTC)
	}
	return start
}

func main() {
	db, _ := sql.Open("postgres", "postgres://...")
	uc := usecase.SubmitLeave{
		Employees: drivers.PostgresEmployeeRepo{DB: db},
		Leaves:    drivers.PostgresLeaveRepo{DB: db},
		Mailer:    drivers.SMTPMailer{},
		Clock:     sysClock{},
		YearStart: fiscalYearStart,
	}
	http.Handle("/leave-requests", adapters.SubmitHandler{UC: uc})
	http.ListenAndServe(":8080", nil)
}

// usecase/submit_leave.go
package usecase

import (
	"errors"
	"time"

	"github.com/ohagi/clean-architecture-examples/good/domain"
)

var (
	ErrNotEligible = errors.New("employee not eligible")
)

type Clock interface{ Now() time.Time }

type EmployeeRepo interface {
	FindByID(id string) (domain.Employee, error)
}

type LeaveRepo interface {
	CountThisFiscalYear(employeeID string, fiscalYearStart time.Time) (int, error)
	Create(req *domain.LeaveRequest) error
}

type Mailer interface {
	NotifyManagerNewRequest(requestID string) error
}

type SubmitLeave struct {
	Employees EmployeeRepo
	Leaves    LeaveRepo
	Mailer    Mailer
	Clock     Clock
	YearStart func(now time.Time) time.Time // 会計年度開始日の計算
}

type SubmitInput struct {
	EmployeeID string
	Reason     string
	From       time.Time
	To         time.Time
}

type SubmitOutput struct {
	ID     string
	Status domain.LeaveStatus
}

func (uc SubmitLeave) Exec(in SubmitInput) (SubmitOutput, error) {
	now := uc.Clock.Now()
	emp, err := uc.Employees.FindByID(in.EmployeeID)
	if err != nil {
		return SubmitOutput{}, err
	}
	count, err := uc.Leaves.CountThisFiscalYear(in.EmployeeID, uc.YearStart(now))
	if err != nil {
		return SubmitOutput{}, err
	}
	if !domain.CanSubmit(emp, count, now) {
		return SubmitOutput{}, ErrNotEligible
	}
	req := &domain.LeaveRequest{
		EmployeeID: in.EmployeeID,
		Reason:     in.Reason,
		From:       in.From,
		To:         in.To,
		Status:     domain.StatusPending,
		CreatedAt:  now,
	}
	if err := uc.Leaves.Create(req); err != nil {
		return SubmitOutput{}, err
	}
	_ = uc.Mailer.NotifyManagerNewRequest(req.ID) // 失敗は致命にしない方針
	return SubmitOutput{ID: req.ID, Status: req.Status}, nil
}

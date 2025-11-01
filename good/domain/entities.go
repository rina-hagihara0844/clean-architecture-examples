// domain/entities.go
package domain

import "time"

type LeaveStatus string

const (
	StatusPending  LeaveStatus = "PENDING"
	StatusApproved LeaveStatus = "APPROVED"
	StatusRejected LeaveStatus = "REJECTED"
	StatusReturned LeaveStatus = "RETURNED" // 差し戻し
)

type Employee struct {
	ID       string
	HireDate time.Time
}

type LeaveRequest struct {
	ID         string
	EmployeeID string
	Reason     string
	From       time.Time
	To         time.Time
	Status     LeaveStatus
	CreatedAt  time.Time
}

// Policy（ルール）: 半年以上、年度5回未満
func CanSubmit(e Employee, submittedCountThisFiscal int, now time.Time) bool {
	sixMonths := now.AddDate(0, -6, 0)
	if e.HireDate.After(sixMonths) {
		return false
	}
	return submittedCountThisFiscal < 5
}

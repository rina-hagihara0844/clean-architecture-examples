package domain

import "time"

type Employee struct {
	ID       string
	Name     string
	HireDate time.Time
}

func (e Employee) IsEligible(now time.Time) bool {
	return e.HireDate.AddDate(0, 6, 0).Before(now)
}

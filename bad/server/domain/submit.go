package domain

import (
	"fmt"
	"time"
)

type Submit struct {
	EmployeeID string
	Reason     string
	From       string
	To         string
}

func (s Submit) Validate(emp Employee, submitCnt int) error {
	// 入社6か月未満なら不適格
	if !emp.IsEligible(time.Now()) {
		return fmt.Errorf("not eligible before 6 months after hire")
	}
	// 年5回まで
	if submitCnt >= 5 {
		return fmt.Errorf("over submit limit by year")
	}
	return nil
}

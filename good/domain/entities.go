package domain

// Domain層（ビジネスルール・ビジネス概念の定義）
// --------------------------------------------------------
// この層の責務：
// - 業務上の「ビジネス概念（Entity）」を表現する（Employee, LeaveRequestなど）
// - 業務に関するルール（Policy / Business Rule）を定義する（例：申請可能条件など）
// - 外部技術（DB, HTTP, フレームワーク等）に依存しない純粋なロジックを保つ
// --------------------------------------------------------
// ※この層はアプリケーションの中核。
//   外部の技術的要素から最も独立しており、どの環境でも再利用できる。
// --------------------------------------------------------

import "time"

type LeaveStatus string

const (
	StatusPending  LeaveStatus = "PENDING"
	StatusApproved LeaveStatus = "APPROVED"
	StatusRejected LeaveStatus = "REJECTED"
	StatusReturned LeaveStatus = "RETURNED" // 差し戻し
)

// Employee（従業員）
// ドメインオブジェクト：システム内で従業員を表す純粋なモデル
type Employee struct {
	ID       string
	HireDate time.Time
}

// LeaveRequest（休暇申請）
// ドメインオブジェクト：休暇申請のビジネス上の状態を表す純粋なモデル
type LeaveRequest struct {
	ID         string
	EmployeeID string
	Reason     string
	From       time.Time
	To         time.Time
	Status     LeaveStatus
	CreatedAt  time.Time
}

// ビジネスルール
// 半年以上勤務している & 年度内5回未満なら申請可能
func CanSubmit(e Employee, submittedCountThisFiscal int, now time.Time) bool {
	if e.HireDate.AddDate(0, 6, 0).After(now) {
		return false
	}
	return submittedCountThisFiscal < 5
}

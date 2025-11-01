// usecase/ports.go

// Repositoryや外部連携を抽象化するインターフェース群
// --------------------------------------------------------
// Domain層のルールを実行するために必要なデータや通知を取得・送信するための「契約」。
// 「どんなデータが必要で」「どんな外部処理を呼び出すか」を定義する。
// 必要なデータが取得でき、必要な外部処理が実行できれば実装は問わない。
// 実際の実装（DBアクセスやメール送信）は Infrastructure（Drivers）層が担当。
// --------------------------------------------------------

package usecase

import (
	"time"

	"github.com/ohagi/clean-architecture-examples/good/domain"
)

// Repositoryや外部連携を抽象化するインターフェース群
// --------------------------------------------------------
// Domain層のルールを実行するために必要なデータや通知を取得・送信するための「契約」。
// 「どんなデータが必要で」「どんな外部処理を呼び出すか」を定義する。
// 必要なデータが取得でき、必要な外部処理が実行できれば実装は問わない。
// 実際の実装（DBアクセスやメール送信）は Infrastructure層が担当。
// --------------------------------------------------------
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

package drivers

// Framework & Drivers層（インフラストラクチャ層）
// --------------------------------------------------------
// この層の責務：
// - 外部システム（例：データベース、メールサーバ、外部APIなど）との接続・入出力を扱う
// - Domain / UseCase 層が定義したインターフェース（Repository, Mailerなど）を具象実装する
// - ログ記録・設定読み込み・通信処理など、技術的な詳細を担当する
// --------------------------------------------------------
// この層では業務ロジックを一切持たず、あくまで「実装の詳細」に徹する。
// --------------------------------------------------------
// ※アプリケーションの最も外側に位置し、外部の技術要素を内側の層から隔離する役割を持つ。
//   内側（Domain / UseCase）はこの層に依存しない。
// --------------------------------------------------------

import (
	"database/sql"
	"time"

	"github.com/ohagi/clean-architecture-examples/good/domain"
)

// PostgresEmployeeRepo は従業員情報を PostgreSQL から取得するリポジトリ。
// Domain層の EmployeeRepository インターフェースを満たす。
// 「Domain層の EmployeeRepository インターフェース」というのは、ドメインやユースケースが外部に対して「こういうデータが欲しい」という依頼の窓口（契約）
type PostgresEmployeeRepo struct{ DB *sql.DB }

// FindByID は従業員IDで Employee を検索する。
// 純粋にDBからデータを取得するのみで、業務ルールは扱わない。
func (r PostgresEmployeeRepo) FindByID(id string) (domain.Employee, error) {
	var e domain.Employee
	return e, r.DB.QueryRow(`SELECT id, hire_date FROM employees WHERE id=$1`, id).Scan(&e.ID, &e.HireDate)
}

// PostgresLeaveRepo は休暇申請データを PostgreSQL に保存・取得するリポジトリ。
// Domain層の LeaveRepository インターフェースを満たす。
type PostgresLeaveRepo struct{ DB *sql.DB }

// CountThisFiscalYear は年度内の申請回数をDBからカウントする。
// ビジネス条件（年度開始日など）はUseCaseから与えられる。
func (r PostgresLeaveRepo) CountThisFiscalYear(empID string, start time.Time) (int, error) {
	var c int
	return c, r.DB.QueryRow(
		`SELECT COUNT(*) FROM leave_requests WHERE employee_id=$1 AND created_at >= $2`,
		empID, start).Scan(&c)
}

// Create は新しい休暇申請をDBに登録する。
// 登録時の業務ルール（件数制限・勤務期間チェック等）はUseCase/Domain側で担保される。
func (r PostgresLeaveRepo) Create(req *domain.LeaveRequest) error {
	return r.DB.QueryRow(
		`INSERT INTO leave_requests(employee_id,reason,from_date,to_date,status,created_at)
		 VALUES($1,$2,$3,$4,$5,$6) RETURNING id`,
		req.EmployeeID, req.Reason, req.From, req.To, req.Status, req.CreatedAt,
	).Scan(&req.ID)
}

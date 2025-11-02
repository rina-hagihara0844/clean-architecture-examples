package usecase

//
// UseCase層（アプリケーションロジック）
// --------------------------------------------------------
// この層の責務：
// - Domain層のビジネスルールを「いつ・どの順に」実行するかを定義する
// - 業務フローの制御（例：休暇申請の一連の流れ）
// - 外部との入出力はインターフェース経由で行い、具象技術(DB, HTTP, メール等)には依存しない
// --------------------------------------------------------
// この層は「ビジネスの手続き（シナリオ）」を定義する層。
// Domain層が「ルール」を、UseCase層が「手順」を担う。
// --------------------------------------------------------

import (
	"errors"
	"time"

	"github.com/ohagi/clean-architecture-examples/good/domain"
)

var (
	ErrNotEligible = errors.New("employee not eligible")
)

// SubmitLeave：休暇申請ユースケースの実行構造体
// --------------------------------------------------------
// - 休暇申請処理の全体フローを司る
// - 外部リソースには依存せず、インターフェース経由で操作する
// --------------------------------------------------------
type SubmitLeave struct {
	EmployeesRepo EmployeeRepo
	LeavesRepo    LeaveRepo
	Mailer        Mailer
	Clock         Clock
	YearStart     func(now time.Time) time.Time // 会計年度開始日の計算
}

// SubmitInput / SubmitOutput
// --------------------------------------------------------
// - UseCaseに入力されるデータと、出力される結果を表すDTO（アプリケーション層用）
// - Adapter層（例：HTTPハンドラ）がこれらを使ってデータを受け渡す
// --------------------------------------------------------
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

// Exec：休暇申請ユースケースの実行
// --------------------------------------------------------
// 処理フロー：
// 1. 従業員情報の取得
// 2. 年度内の申請回数の取得
// 3. ドメインルールによる申請可否判定
// 4. 申請データの生成と保存
// 5. 管理者への通知（失敗は致命エラーにしない）
// --------------------------------------------------------
func (uc SubmitLeave) Submit(in SubmitInput) (SubmitOutput, error) {
	now := uc.Clock.Now()

	// 1. 従業員情報の取得
	emp, err := uc.EmployeesRepo.FindByID(in.EmployeeID)
	if err != nil {
		return SubmitOutput{}, err
	}

	// 2. 年度内の申請回数の取得
	count, err := uc.LeavesRepo.CountThisFiscalYear(in.EmployeeID, uc.YearStart(now))
	if err != nil {
		return SubmitOutput{}, err
	}

	// 3. ドメインルールによる申請可否判定
	if !domain.CanSubmit(emp, count, now) {
		return SubmitOutput{}, ErrNotEligible
	}
	// 休暇申請データの生成
	req := &domain.LeaveRequest{
		EmployeeID: in.EmployeeID,
		Reason:     in.Reason,
		From:       in.From,
		To:         in.To,
		Status:     domain.StatusPending,
		CreatedAt:  now,
	}
	// 4. 申請データの生成と保存
	if err := uc.LeavesRepo.Create(req); err != nil {
		return SubmitOutput{}, err
	}
	if err := uc.Mailer.NotifyManagerNewRequest(req.ID); err != nil {
		return SubmitOutput{}, err
	}
	return SubmitOutput{ID: req.ID, Status: req.Status}, nil
}

package adapters

// Interface Adapter層（インターフェースアダプター）
// --------------------------------------------------------
// この層の責務：
// - 外部入出力（例：HTTP, CLI, GraphQL など）とアプリケーション内部との橋渡しを行う
// - HTTPリクエストをアプリケーションで扱える形式（DTO）へ変換
// - UseCaseを呼び出して業務処理を実行
// - UseCaseの出力をHTTPレスポンス形式に変換して返す
// --------------------------------------------------------
// この層の目的は「データの変換とルーティング」に限定される。
// ビジネスロジック（ルールや判断）はここに書かない。
// --------------------------------------------------------
// ※Framework層（HTTPサーバなど）とUseCase層の間を仲介する「翻訳者」。
//   外の世界（HTTPなど）を内側（UseCase）に安全に渡す役割を担う。
// --------------------------------------------------------

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ohagi/clean-architecture-examples/good/domain"
	"github.com/ohagi/clean-architecture-examples/good/usecase"
)

// Usecaseのインターフェースを持つ
// これにより、Handler は具体的な業務ロジックを知らずに UseCase を呼び出せる
type SubmitHandler struct{ UC usecase.SubmitLeave }

func (h SubmitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// HTTPリクエストボディをGo構造体にパースするためのDTO
	var body struct {
		EmployeeID string `json:"employeeId"`
		Reason     string `json:"reason"`
		From       string `json:"from"`
		To         string `json:"to"`
	}
	// JSONデコード
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad json", 400)
		return
	}
	// 日付パース
	from, err1 := time.Parse("2006-01-02", body.From)
	to, err2 := time.Parse("2006-01-02", body.To)
	if err1 != nil || err2 != nil {
		http.Error(w, "bad date", 400)
		return
	}
	// UseCaseの呼び出し
	out, err := h.UC.Submit(usecase.SubmitInput{
		EmployeeID: body.EmployeeID, Reason: body.Reason, From: from, To: to,
	})
	// エラーハンドリング
	if err != nil {
		status := 400
		if err == usecase.ErrNotEligible {
			status = 403
		}
		http.Error(w, err.Error(), status)
		return
	}
	// 成功レスポンスの返却
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(struct {
		ID     string             `json:"id"`
		Status domain.LeaveStatus `json:"status"`
	}{out.ID, out.Status})
}

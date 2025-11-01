// server_bad.go  — SRP違反：HandlerがDomain/UseCase/Infraを全部抱える
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

// HTTPリクエストボディ
// → Interface Adapter層（HTTPの入力DTO）
type SubmitReq struct {
	EmployeeID string `json:"employeeId"`
	Reason     string `json:"reason"`
	From       string `json:"from"` // "2025-11-01"
	To         string `json:"to"`
}

func main() {
	// → Framework & Driver層（Infraの具体実装）
	db, _ := sql.Open("postgres", "postgres://...")
	defer db.Close()

	// HTTPハンドラ登録
	// → Interface Adapter層
	http.HandleFunc("/leave-requests", func(w http.ResponseWriter, r *http.Request) {
		// HTTPリクエストボディ（JSON）を Go の構造体 SubmitReq に変換
		//  → Interface Adapter層
		var req SubmitReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", 400)
			return
		}

		// 正規表現チェック
		// → Interface Adapter層
		reISO := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
		if req.EmployeeID == "" || !reISO.MatchString(req.From) || !reISO.MatchString(req.To) {
			http.Error(w, "validation error", 400)
			return
		}

		// DBアクセス直書き
		// → Framework & Driver層
		var hireDate time.Time
		if err := db.QueryRow(`SELECT hire_date FROM employees WHERE id=$1`, req.EmployeeID).Scan(&hireDate); err != nil {
			http.Error(w, "employee not found", 404)
			return
		}
		// 入社6か月以上かチェック
		// → Domain層
		if hireDate.AddDate(0, 6, 0).Before(time.Now()) {
			http.Error(w, "not eligible (<6 months)", 400)
			return
		}

		// ----------------------------
		// 年度内の申請回数チェック
		// → usecase層
		// ----------------------------
		yearStart := time.Date(time.Now().Year(), 4, 1, 0, 0, 0, 0, time.UTC) // 例：会計年度=4月開始
		var cnt int
		// DBアクセス直書き
		// 対象の従業員の今年度の申請回数をカウント
		// → Framework & Driver層
		if err := db.QueryRow(`SELECT COUNT(*) FROM leave_requests WHERE employee_id=$1 AND created_at>= $2 AND status IN ('PENDING','APPROVED','REJECTED','RETURNED')`,
			req.EmployeeID, yearStart).Scan(&cnt); err != nil {
			http.Error(w, "db error", 500)
			return
		}
		// 年5回まで
		// → Domain層
		if cnt >= 5 {
			http.Error(w, "over yearly limit", 400)
			return
		}

		// ----------------------------
		// UseCase層
		// - 休暇申請を登録
		// ----------------------------
		var id string
		// DBアクセス直書き
		// → Framework & Driver層
		if err := db.QueryRow(
			`INSERT INTO leave_requests(employee_id,reason,from_date,to_date,status,created_at)
			 VALUES($1,$2,$3,$4,'PENDING',NOW()) RETURNING id`,
			req.EmployeeID, req.Reason, req.From, req.To).Scan(&id); err != nil {
			http.Error(w, "db error", 500)
			return
		}

		// メール送信直書き
		// → Framework & Driver層
		_ = sendSMTP("manager@example.com", "New request "+id)

		// HTTPレスポンスボディ（JSON）を返す
		// → Interface Adapter層
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": id, "status": "PENDING",
		})
	})

	// サーバ起動
	// → Framework & Driver層
	http.ListenAndServe(":8080", nil)
}

// メール送信（ダミー実装）
// → Framework & Driver層
func sendSMTP(to, body string) error {
	fmt.Println("to:", to)
	fmt.Println("body:", body)
	return nil
}

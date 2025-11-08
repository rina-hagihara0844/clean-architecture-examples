package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	ad_db "example.com/server/adapter/db"
	"example.com/server/adapter/mail"
	"example.com/server/domain"
)

// HTTPリクエストボディ
// HTTPの入力DTO
type SubmitReq struct {
	EmployeeID string `json:"employeeId"`
	Reason     string `json:"reason"`
	From       string `json:"from"` // "2025-11-01"
	To         string `json:"to"`
}

type SubmitRes struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (r *SubmitReq) toEntity() (domain.Submit, error) {
	reISO := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if r.EmployeeID == "" || !reISO.MatchString(r.From) || !reISO.MatchString(r.To) {
		return domain.Submit{}, fmt.Errorf("validation error")
	}
	return domain.Submit{
		EmployeeID: r.EmployeeID,
		Reason:     r.Reason,
		From:       r.From,
		To:         r.To,
	}, nil
}

func leaveRequestParser(r *http.Request) (SubmitReq, error) {
	var req SubmitReq
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func LeaveRequest(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// HTTPリクエストボディ（JSON）を Go の構造体 SubmitReq に変換
	req, err := leaveRequestParser(r)
	if err != nil {
		http.Error(w, "bad json", 400)
		return
	}

	// フォーマットチェック(toEntity)
	submit, err := req.toEntity()
	if err != nil {
		http.Error(w, "validation error", 400)
		return
	}

	// 従業員情報取得
	emp, err := ad_db.GetEmployeeByID(db, submit.EmployeeID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "employee not found", 404)
			return
		}
		http.Error(w, "db error", 500)
		return
	}

	// 年度内の申請回数チェック
	cnt, err := ad_db.CountSubmitThisFiscalYear(db, submit.EmployeeID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "employee not found", 404)
			return
		}
		http.Error(w, "db error", 500)
		return
	}

	// 休暇申請の妥当性チェック
	if err := submit.Validate(emp, cnt); err != nil {
		http.Error(w, "validation error: "+err.Error(), 400)
		return
	}

	// 休暇申請を登録
	id, err := ad_db.RegisterLeaveRequest(db, submit)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "employee not found", 404)
			return
		}
		// 409もやりたい
		http.Error(w, "db error", 500)
		return
	}

	// メール送信
	mail.SendSMTP("manager@example.com", "New request "+id)

	// HTTPレスポンスボディ（JSON）を返す
	w.Header().Set("Content-Type", "application/json")
	res := SubmitRes{ID: id, Status: "PENDING"}
	json.NewEncoder(w).Encode(res)
	// このjsonを返す
}

package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/ohagi/clean-architecture-examples/good/adapters"
	"github.com/ohagi/clean-architecture-examples/good/drivers"
	"github.com/ohagi/clean-architecture-examples/good/usecase"
)

type sysClock struct{}

func (sysClock) Now() time.Time { return time.Now().UTC() }
func fiscalYearStart(now time.Time) time.Time {
	year := now.Year()
	start := time.Date(year, 4, 1, 0, 0, 0, 0, time.UTC)
	if now.Before(start) {
		start = time.Date(year-1, 4, 1, 0, 0, 0, 0, time.UTC)
	}
	return start
}

func main() {
	// DB接続の初期化
	db, _ := sql.Open("postgres", "postgres://...")
	// 依存性の注入
	// UseCaseはインターフェイスに依存するので、ここで具体実装を差し込む
	uc := usecase.SubmitLeave{
		EmployeesRepo: drivers.PostgresEmployeeRepo{DB: db},
		LeavesRepo:    drivers.PostgresLeaveRepo{DB: db},
		Mailer:        drivers.SMTPMailer{},
		Clock:         sysClock{},
		YearStart:     fiscalYearStart,
	}
	// HTTPハンドラの登録
	// HandlerにはUseCaseを注入して利用する
	http.Handle("/leave-requests", adapters.SubmitHandler{UC: uc})
	// HTTPサーバ起動
	http.ListenAndServe(":8080", nil)
}

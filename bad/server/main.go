// server_bad.go  — SRP違反：HandlerがDomain/UseCase/Infraを全部抱える
package main

import (
	"database/sql"
	"net/http"

	"example.com/server/handler"
)

func main() {
	// DB接続初期化
	db, _ := sql.Open("postgres", "postgres://...")
	defer db.Close()

	// HTTPハンドラ登録
	http.HandleFunc("/leave-requests", func(w http.ResponseWriter, r *http.Request) {
		handler.LeaveRequest(db, w, r)
	})

	http.ListenAndServe(":8080", nil)
}

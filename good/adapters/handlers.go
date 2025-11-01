// adapters/http_handlers.go
package adapters

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ohagi/clean-architecture-examples/good/domain"
	"github.com/ohagi/clean-architecture-examples/good/usecase"
)

type SubmitHandler struct{ UC usecase.SubmitLeave }

func (h SubmitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		EmployeeID string `json:"employeeId"`
		Reason     string `json:"reason"`
		From       string `json:"from"`
		To         string `json:"to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad json", 400)
		return
	}
	from, err1 := time.Parse("2006-01-02", body.From)
	to, err2 := time.Parse("2006-01-02", body.To)
	if err1 != nil || err2 != nil {
		http.Error(w, "bad date", 400)
		return
	}
	out, err := h.UC.Exec(usecase.SubmitInput{
		EmployeeID: body.EmployeeID, Reason: body.Reason, From: from, To: to,
	})
	if err != nil {
		status := 400
		if err == usecase.ErrNotEligible {
			status = 403
		}
		http.Error(w, err.Error(), status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(struct {
		ID     string             `json:"id"`
		Status domain.LeaveStatus `json:"status"`
	}{out.ID, out.Status})
}

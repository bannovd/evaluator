package application

import (
	"bytes"
	"net/http"

	"github.com/bannovd/evaluator/models"
)

// HealthHandler check is alive
func (app *Application) HealthHandler(w http.ResponseWriter, r *http.Request) {
	sum := models.GetHashSum()
	if !bytes.Equal(sum, app.hashSum) {
		_ = app.logger.Log("msg", "New Configuration")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

package application

import (
	"errors"
	"net/http"

	"github.com/bannovd/evaluator/models"
)

// CatchHandler check is alive
func (app *Application) CatchHandler(w http.ResponseWriter, r *http.Request) {
	hit, err := ParseQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.svc.Repository.SaveHit(*hit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// ParseQuery func
func ParseQuery(r *http.Request) (*models.Hit, error) {
	query := r.URL.Query()

	val := query.Get("val")
	if val == "" {
		return nil, errors.New("Parameter 'val' is missing.")
	}

	t := query.Get("t")
	if t == "" {
		return nil, errors.New("Parameter 't' is missing.")
	}

	switch t {
	case "pageview", "event":
		hit := models.Hit{
			Type:  t,
			Value: val}
		return &hit, nil
	default:
		return nil, errors.New("Unknown parameter 't'.")
	}
}

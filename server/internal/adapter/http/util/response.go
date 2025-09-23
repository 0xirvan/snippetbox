package util

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func WriteValidationError(w http.ResponseWriter, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		out := make(map[string]string)
		for _, e := range errs {
			out[e.Field()] = e.Tag()
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(out)
		return
	}
	http.Error(w, err.Error(), http.StatusBadRequest)
}

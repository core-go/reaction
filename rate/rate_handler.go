package rate

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func NewRateHandler(
	service RateService,
	authorIndex int,
	idIndex int,
	max int,
) Handler {
	return Handler{
		service:     service,
		authorIndex: authorIndex,
		idIndex:     idIndex,
		max:         max,
	}
}

type Handler struct {
	service     RateService
	authorIndex int
	idIndex     int
	max         int
}

func (h *Handler) Rate(w http.ResponseWriter, r *http.Request) {
	var rate Request
	er1 := Decode(w, r, &rate)
	author := GetRequiredParam(w, r, h.authorIndex)
	id := GetRequiredParam(w, r, h.idIndex)

	if er1 == nil {
		errors := Validate(r.Context(), rate, h.max)
		if len(errors) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(errors)
			return
		}
		result, er3 := h.service.Rate(r.Context(), id, author, rate)
		if er3 != nil {
			http.Error(w, er3.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
	}
}

func GetRequiredParam(w http.ResponseWriter, r *http.Request, options ...int) string {
	p := GetParam(r, options...)
	if len(p) == 0 {
		http.Error(w, "parameter is required", http.StatusBadRequest)
		return ""
	}
	return p
}

func GetParam(r *http.Request, options ...int) string {
	offset := 0
	if len(options) > 0 && options[0] > 0 {
		offset = options[0]
	}
	s := r.URL.Path
	params := strings.Split(s, "/")
	i := len(params) - 1 - offset
	if i >= 0 {
		return params[i]
	} else {
		return ""
	}
}

func Decode(w http.ResponseWriter, r *http.Request, obj interface{}, options ...func(context.Context, interface{}) (interface{}, error)) error {
	er1 := json.NewDecoder(r.Body).Decode(obj)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return er1
	}
	if len(options) > 0 && options[0] != nil {
		_, er2 := options[0](r.Context(), obj)
		if er2 != nil {
			http.Error(w, er2.Error(), http.StatusInternalServerError)
		}
		return er2
	}
	return nil
}

func Validate(ctx context.Context, rate Request, max int) []ErrorMessage {
	var errors []ErrorMessage
	if rate.Rate > max {
		errors = append(errors, ErrorMessage{
			Field: "rate",
			Code:  "max",
			Param: strconv.Itoa(max),
		})
	}
	return errors
}

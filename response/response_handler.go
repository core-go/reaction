package response

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ResponseHandler interface {
	Load(w http.ResponseWriter, r *http.Request)
	Response(w http.ResponseWriter, r *http.Request)
}

func NewResponseHandler(
	service ResponseService,
	idField string,
	authorField string,
) ResponseHandler {
	return &responseHandler{service: service, idField: idField, authorField: authorField}
}

type responseHandler struct {
	service     ResponseService
	idField     string
	authorField string
}

func (h *responseHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)[h.idField]
	author := mux.Vars(r)[h.authorField]
	if len(id) > 0 && len(author) > 0 {
		res, err := h.service.Load(r.Context(), id, author)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(res)
		return
	}
	http.Error(w, "parameter is required", http.StatusBadRequest)
	return
}

func (h *responseHandler) Response(w http.ResponseWriter, r *http.Request) {
	var response Response
	var t = time.Now()
	response.Time = &t
	er1 := Decode(w, r, &response)
	if er1 != nil {
		return
	}
	response.Id = mux.Vars(r)[h.idField]
	response.Author = mux.Vars(r)[h.authorField]

	if len(response.Id) == 0 || len(response.Author) == 0 {
		http.Error(w, "parameter is required", http.StatusBadRequest)
		return
	}

	res, er3 := h.service.Response(r.Context(), &response)
	if er3 != nil {
		http.Error(w, er3.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(res)
	return

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

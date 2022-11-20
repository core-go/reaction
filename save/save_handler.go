package save

import (
	"encoding/json"
	"net/http"
	"strings"
)

func NewSaveHandler(service SaveService, itemIndex int, idIndex int) SaveHandler {
	return SaveHandler{service: service, itemIndex: itemIndex, idIndex: idIndex}
}

type SaveHandler struct {
	service SaveService
	itemIndex int
	idIndex int
}

func (h *SaveHandler) Save(w http.ResponseWriter, r *http.Request) {
	item := GetRequiredParam(w, r, h.itemIndex) // 0
	id := GetRequiredParam(w, r, h.idIndex) // 1
	if len(id) > 0 && len(item) > 0 {
		result, err := h.service.Save(r.Context(), id, item)
		if err != nil {
			http.Error(w,err.Error(),500)
		return;
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
		return;
	}
}

func (h *SaveHandler) Remove(w http.ResponseWriter, r *http.Request) {
	item := GetRequiredParam(w, r,h.itemIndex) // 0
	id := GetRequiredParam(w, r, h.idIndex) // 1
	if len(id) > 0 && len(item) > 0 {
		result, err := h.service.Remove(r.Context(), id, item)
		if err != nil {
			http.Error(w,err.Error(),500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
		return
	}
}

func (h *SaveHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := GetRequiredParam(w, r)
	if len(id) > 0 {
		result, err := h.service.Load(r.Context(), id)
		if err != nil {
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
		return
	}
}

func GetRequiredParam(w http.ResponseWriter,r *http.Request, options ...int) string {
	p := GetParam(r, options...)
	if len(p) == 0 {
		http.Error(w, "parameter is required", http.StatusBadRequest)
		return ""
	}
	return p
}
func GetParam(r *http.Request, options... int) string {
	offset := 0
	if len(options) > 0 && options[0] > 0 {
		offset = options[0]
	}
	s := r.URL.Path
	params := strings.Split(s, "/")
	i := len(params)-1-offset
	if i >= 0 {
		return params[i]
	} else {
		return ""
	}
}

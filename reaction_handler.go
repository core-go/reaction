package reaction

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func NewReactionHandler(
	service ReactionService,
	userIdIndex int,
	authorIndex int,
	idIndex int,
) ReactionHandler {
	return ReactionHandler{
		service:     service,
		userIdIndex: userIdIndex,
		authorIndex: authorIndex,
		idIndex:     idIndex,
	}
}

type ReactionHandler struct {
	service     ReactionService
	userIdIndex int
	authorIndex int
	idIndex     int
}

func (h *ReactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reaction Reaction
	var t = time.Now()
	er1 := Decode(w, r, &reaction)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
	}
	reaction.UserId = GetRequiredParam(w, r, h.userIdIndex)
	reaction.Author = GetRequiredParam(w, r, h.authorIndex)
	reaction.Id = GetRequiredParam(w, r, h.idIndex)

	reaction.Time = &t
	if reaction.Type == 0 {
		reaction.Type = 1
	}

	result, er3 := h.service.Insert(r.Context(), &reaction)
	if er3 != nil {
		http.Error(w, er3.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(result)
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
func (h *ReactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var reaction Reaction
	reaction.UserId = GetRequiredParam(w, r, h.userIdIndex)
	reaction.Author = GetRequiredParam(w, r, h.authorIndex)
	reaction.Id = GetRequiredParam(w, r, h.idIndex)
	result, err := h.service.Delete(r.Context(), &reaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(result)
	return
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
func GetRequiredParam(w http.ResponseWriter, r *http.Request, options ...int) string {
	p := GetParam(r, options...)
	if len(p) == 0 {
		http.Error(w, "parameter is required", http.StatusBadRequest)
		return ""
	}
	return p
}

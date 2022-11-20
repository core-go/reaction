package userreaction

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func NewUserReactionHandler(
	service UserReactionService,
	idField string,
	authorField string,
	reactionField string,
) UserReactionHandler {
	return UserReactionHandler{service: service, idField: idField, authorField: authorField, reactionField: reactionField}
}

type UserReactionHandler struct {
	service       UserReactionService
	idField       string
	authorField   string
	reactionField string
}

func (h *UserReactionHandler) CheckReact(w http.ResponseWriter, r *http.Request) {
	//id := GetRequiredParam(w, r, 1)
	//author := GetRequiredParam(w, r,0)
	id := mux.Vars(r)[h.idField]
	author := mux.Vars(r)[h.authorField]
	if len(id) > 0 && len(author) > 0 {
		res, err := h.service.CheckReaction(r.Context(), id, author)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(res)
		return
	}
}

func (h *UserReactionHandler) Unreact(w http.ResponseWriter, r *http.Request) {
	//id := GetRequiredParam(w, r, 2)
	//author := GetRequiredParam(w, r, 1)
	//reaction := GetRequiredParam(w, r, 0)
	id := mux.Vars(r)[h.idField]
	author := mux.Vars(r)[h.authorField]
	reaction := mux.Vars(r)[h.reactionField]
	if len(id) > 0 && len(author) > 0 && len(reaction) > 0 {
		res, err := h.service.Unreact(r.Context(), id, author, reaction)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(res)
		return
	}
}

func (h *UserReactionHandler) React(w http.ResponseWriter, r *http.Request) {
	//id := GetRequiredParam(w, r, 2)
	//author := GetRequiredParam(w, r, 1)
	//reaction := GetRequiredParam(w, r, 0)
	id := mux.Vars(r)[h.idField]
	author := mux.Vars(r)[h.authorField]
	reaction := mux.Vars(r)[h.reactionField]
	if len(id) > 0 && len(author) > 0 && len(reaction) > 0 {
		res, err := h.service.React(r.Context(), id, author, reaction)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(res)
		return
	}
	http.Error(w, "parameter is required", http.StatusInternalServerError)
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

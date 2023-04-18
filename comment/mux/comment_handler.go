package mux

import (
	"context"
	"encoding/json"
	"github.com/core-go/reaction/comment"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func NewCommentHandler(
	service comment.CommentService,
	commentIdField string,
	idField string,
	authorField string,
	userIdField string,
) CommentHandler {
	return CommentHandler{service: service, commentIdField: commentIdField, idField: idField, authorField: authorField, userIdField: userIdField}
}

type CommentHandler struct {
	service        comment.CommentService
	commentIdField string
	idField        string
	authorField    string
	userIdField    string
}

func (h *CommentHandler) Load(w http.ResponseWriter, r *http.Request) {
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
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var comment comment.Request
	er1 := Decode(w, r, &comment)
	if er1 != nil {
		return
	}
	id := mux.Vars(r)[h.idField]
	commentId := uuid.New().String()
	author := mux.Vars(r)[h.authorField]
	userId := mux.Vars(r)[h.userIdField]
	if len(author) > 0 && len(id) > 0 {
		res, er3 := h.service.Create(r.Context(), id, commentId, userId, author, comment)
		if er3 != nil {
			http.Error(w, er3.Error(), 500)
			return
		}
		if res <= 0 {
			http.Error(w, "no records affected", http.StatusInternalServerError)
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

func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	var comment comment.Request
	er1 := Decode(w, r, &comment)
	if er1 != nil {
		return
	}
	id := mux.Vars(r)[h.idField]
	commentId := mux.Vars(r)[h.commentIdField]
	author := mux.Vars(r)[h.authorField]
	userid := mux.Vars(r)[h.userIdField]
	if len(commentId) <= 0 || len(author) <= 0 || len(id) <= 0 {
		http.Error(w, "paramerter is required", http.StatusBadRequest)
		return
	}
	res, er3 := h.service.Update(r.Context(), id, commentId, userid, author, comment)
	if er3 != nil {
		http.Error(w, er3.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(res)
	return
}

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var commentId, id, author string
	commentId = mux.Vars(r)[h.commentIdField]
	author = mux.Vars(r)[h.authorField]
	id = mux.Vars(r)[h.idField]
	if len(commentId) > 0 && len(author) > 0 && len(id) > 0 {
		res, err := h.service.Delete(r.Context(), id, commentId, author)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(res)
		return
	}
	http.Error(w, "paramerter is required", http.StatusBadRequest)
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

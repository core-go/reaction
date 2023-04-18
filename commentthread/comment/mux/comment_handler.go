package mux

import (
	"context"
	"encoding/json"
	"github.com/core-go/reaction/commentthread/comment"
	"net/http"

	"github.com/gorilla/mux"
)

type CommentHandler struct {
	service              comment.CommentService
	commentThreadIdField string
	userIdField          string
	authorField          string
	idField              string
	commentIdField       string
	GenerateId           func(context.Context) (string, error)
}

func NewCommentHandler(service comment.CommentService,
	commentThreadIdField string,
	userIdField string,
	authorField string,
	idField string,
	commentIdField string,
	generateId func(context.Context) (string, error)) CommentHandler {
	return CommentHandler{
		service:              service,
		commentThreadIdField: commentThreadIdField,
		userIdField:          userIdField,
		authorField:          authorField,
		idField:              idField,
		commentIdField:       commentIdField,
		GenerateId:           generateId,
	}
}

func (h *CommentHandler) GetReplyComments(w http.ResponseWriter, r *http.Request) {
	obj := make(map[string]string)
	commentThreadId := mux.Vars(r)[h.commentThreadIdField]
	err := Decode(w, r, &obj)
	if err != nil {
		return
	}
	userid := obj[h.userIdField]
	res, err := h.service.GetComments(r.Context(), commentThreadId, &userid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(res)
}

func (h *CommentHandler) Reply(w http.ResponseWriter, r *http.Request) {
	var obj comment.Request
	commentThreadId := mux.Vars(r)[h.commentThreadIdField]
	author := mux.Vars(r)[h.authorField]
	id := mux.Vars(r)[h.idField]
	if len(commentThreadId) == 0 || len(author) == 0 || len(id) == 0 {
		http.Error(w, "parameter is required", http.StatusBadRequest)
		return
	}
	err := Decode(w, r, &obj)
	if err != nil {
		return
	}
	commentId, err := h.GenerateId(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err != nil {
		return
	}
	res, err := h.service.Create(r.Context(), id, commentId, commentThreadId, author, obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res <= 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(res)
}

func (h *CommentHandler) UpdateReply(w http.ResponseWriter, r *http.Request) {
	var obj comment.Request
	commentId := mux.Vars(r)[h.commentIdField]
	author := mux.Vars(r)[h.authorField]
	err := Decode(w, r, &obj)
	if err != nil {
		return
	}
	res, err := h.service.Update(r.Context(), commentId, author, obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res <= 0 {
		if res == -2 {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(res)
}
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	commentId := mux.Vars(r)[h.commentIdField]
	commentThreadId := mux.Vars(r)[h.commentThreadIdField]
	author := mux.Vars(r)[h.authorField]
	if len(commentId) <= 0 || len(commentThreadId) <= 0 || len(author) <= 0 {
		http.Error(w, "parameters is required", http.StatusBadRequest)
		return
	}
	res, err := h.service.Remove(r.Context(), commentId, commentThreadId, author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(res)

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

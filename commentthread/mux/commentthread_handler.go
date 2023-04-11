package mux

import (
	"context"
	"encoding/json"
	commentthread "go-service/internal/reaction/commentthread"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func NewCommentThreadHandler(
	service commentthread.CommentThreadService,
	GenerateId func(ctx context.Context) (string, error),
	commentIdField string,
	authorField string,
	idField string,
) CommentThreadHandler {
	return CommentThreadHandler{
		service:        service,
		generateId:     GenerateId,
		commentIdField: commentIdField,
		idField:        idField,
		authorField:    authorField,
	}
}

type CommentThreadHandler struct {
	service        commentthread.CommentThreadService
	generateId     func(ctx context.Context) (string, error)
	commentIdField string
	authorField    string
	idField        string
}

func (h *CommentThreadHandler) Delete(w http.ResponseWriter, r *http.Request) {
	commentId := mux.Vars(r)[h.commentIdField]
	res, err := h.service.Remove(r.Context(), commentId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(res)
}

func (h *CommentThreadHandler) Comment(w http.ResponseWriter, r *http.Request) {
	var comment commentthread.CommentThread
	t := time.Now()
	er1 := Decode(w, r, &comment)
	if er1 != nil {
		return
	}
	// comment.Id = GetRequiredParam(w, r, 3)
	// comment.Author = GetRequiredParam(w, r, 2)
	comment.Id = mux.Vars(r)[h.idField]
	comment.Author = mux.Vars(r)[h.authorField]
	if len(comment.Id) == 0 || len(comment.Author) == 0 {
		http.Error(w, "parameter is required", http.StatusBadRequest)
		return
	}
	comment.UserId = comment.Author
	comment.CommentId, er1 = h.generateId(r.Context())
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
	}
	comment.Time = t
	result, er3 := h.service.Comment(r.Context(), comment)
	if er3 != nil {
		http.Error(w, er3.Error(), http.StatusInternalServerError)
		return
	}
	if result <= 0 {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(comment.CommentId)
}
func (h *CommentThreadHandler) Update(w http.ResponseWriter, r *http.Request) {
	var comment commentthread.CommentThread
	err := Decode(w, r, &comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment.CommentId = GetRequiredParam(w, r)
	res, err1 := h.service.Update(r.Context(), comment)
	if err1 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(res)

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

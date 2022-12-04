package reaction

import (
	"encoding/json"
	"net/http"
	"strings"
)

type CommentReactionHandler struct {
	service        CommentReactionService
	commentIdIndex int
	authorIndex    int
	userIdIndex    int
}

func NewCommentReactionHandler(service CommentReactionService, commentIdIndex int, authorIdIndex int, userIdIndex int) CommentReactionHandler {
	return CommentReactionHandler{
		service:        service,
		commentIdIndex: commentIdIndex,
		authorIndex:    authorIdIndex,
		userIdIndex:    userIdIndex,
	}
}

func (h *CommentReactionHandler) SetUsetful(w http.ResponseWriter, r *http.Request) {
	commentId := GetRequiredParam(w, r, h.commentIdIndex)
	author := GetRequiredParam(w, r, h.authorIndex)
	userId := GetRequiredParam(w, r, h.userIdIndex)
	if len(commentId) == 0 || len(author) == 0 || len(userId) == 0 {
		http.Error(w, "parameter is required", http.StatusBadRequest)
		return
	}
	res, err := h.service.Save(r.Context(), commentId, author, userId, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(res)
}

func (h *CommentReactionHandler) RemoveUseful(w http.ResponseWriter, r *http.Request) {
	commentId := GetRequiredParam(w, r, h.commentIdIndex)
	author := GetRequiredParam(w, r, h.authorIndex)
	userId := GetRequiredParam(w, r, h.userIdIndex)
	if len(commentId) == 0 || len(author) == 0 || len(userId) == 0 {
		http.Error(w, "parameter is required", http.StatusBadRequest)
		return
	}
	res, err := h.service.Remove(r.Context(), commentId, author, userId)
	if err != nil {
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

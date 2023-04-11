package search

import (
	"context"
	"encoding/json"
	"github.com/core-go/search"
	commentthread2 "go-service/internal/reaction/commentthread"
	"net/http"
)

type SearchResult struct {
	List  []commentthread2.CommentThread `json:"list"`
	Total int64                          `json:"total"`
}

func NewSearchCommentThreadHandler(
	service CommentThreadSearchService,
) *CommentThreadSearchHandler {
	return &CommentThreadSearchHandler{
		service: service,
	}
}

type CommentThreadSearchHandler struct {
	service CommentThreadSearchService
}

func (h *CommentThreadSearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	var filter commentthread2.CommentThreadFilter
	er1 := Decode(w, r, &filter)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}
	if filter.Filter != nil {
		search.RepairFilter(filter.Filter)
	}

	list, total, err := h.service.Search(r.Context(), &filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := SearchResult{
		List:  list,
		Total: total,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(result)
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

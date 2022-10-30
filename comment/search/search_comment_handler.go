package search

import (
	"context"
	"github.com/core-go/search"
	"net/http"
	"reflect"
)

type SearchCommentHandler interface {
	Search(w http.ResponseWriter, r *http.Request)
}

func NewSearchCommentHandler(
	find func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error),
	logError func(context.Context, string, ...map[string]interface{}),
) SearchCommentHandler {
	searchModelType := reflect.TypeOf(CommentFilter{})
	modelType := reflect.TypeOf(Comment{})
	var writeLog func(context.Context, string, string, bool, string) error
	searchHandler := search.NewSearchHandler(find, modelType, searchModelType, logError, writeLog)
	return &searchCommentHandler{SearchHandler: searchHandler}
}

type searchCommentHandler struct {
	*search.SearchHandler
}

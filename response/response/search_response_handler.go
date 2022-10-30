package response

import (
	"context"
	"net/http"
	"reflect"
	"github.com/core-go/search"
)

type SearchResponseHandler interface {
	Search(w http.ResponseWriter, r *http.Request)
}

func NewSearchResponseHandler(
	find func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error),
	logError func(context.Context, string, ...map[string]interface{}),
) SearchResponseHandler {
	searchModelType := reflect.TypeOf(ResponseFilter{})
	modelType := reflect.TypeOf(Response{})
	var writeLog func(context.Context, string, string, bool, string) error
	searchHandler := search.NewSearchHandler(find, modelType, searchModelType, logError, writeLog)
	return &searchResponseHandler{SearchHandler: searchHandler}
}

type searchResponseHandler struct {
	*search.SearchHandler
}

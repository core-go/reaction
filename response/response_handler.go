package response

import (
	"context"
	"net/http"
	"reflect"
	"time"

	sv "github.com/core-go/core"
	"github.com/gorilla/mux"
)

type ResponseHandler interface {
	Load(w http.ResponseWriter, r *http.Request)
	Response(w http.ResponseWriter, r *http.Request)
}

func NewResponseHandler(
	service ResponseService,
	status sv.StatusConfig,
	logError func(context.Context, string, ...map[string]interface{}),
	validate func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error),
	action *sv.ActionConfig,
) ResponseHandler {
	modelType := reflect.TypeOf(Response{})
	params := sv.CreateParams(modelType, &status, logError, validate, action)
	return &responseHandler{service: service, Params: params}
}

type responseHandler struct {
	service ResponseService
	*sv.Params
}

func (h *responseHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	author := mux.Vars(r)["author"]
	if len(id) > 0 {
		res, err := h.service.Load(r.Context(), id, author)
		sv.RespondModel(w, r, res, err, h.Error, nil)
	}
}

func (h *responseHandler) Response(w http.ResponseWriter, r *http.Request) {
	var response Response
	var t = time.Now()
	response.Time = &t
	er1 := sv.Decode(w, r, &response)
	if mux.Vars(r)["id"] != "" {
		response.Id = mux.Vars(r)["id"]
	}
	if mux.Vars(r)["author"] != "" {
		response.Author = mux.Vars(r)["author"]
	}
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &response)
		if !sv.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Create) {
			result, er3 := h.service.Response(r.Context(), &response)
			sv.AfterCreated(w, r, &response, result, er3, h.Status, h.Error, h.Log, h.Resource, h.Action.Create)
		}
	}
}

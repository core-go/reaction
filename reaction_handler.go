package reaction

import (
	"context"
	sv "github.com/core-go/core"
	"net/http"
	"reflect"
	"time"
)

type ReactionHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewReactionHandler(service ReactionService, status sv.StatusConfig, logError func(context.Context, string, ...map[string]interface{}), validate func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error), action *sv.ActionConfig) ReactionHandler {
	modelType := reflect.TypeOf(Reaction{})
	params := sv.CreateParams(modelType, &status, logError, validate, action)
	return &reactionHandler{service: service, Params: params}
}

type reactionHandler struct {
	service ReactionService
	*sv.Params
}

func (h *reactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reaction Reaction
	var t = time.Now()
	er1 := sv.Decode(w, r, &reaction)
	reaction.UserId = sv.GetRequiredParam(w, r, 0)
	reaction.Author = sv.GetRequiredParam(w, r, 2)
	reaction.Id = sv.GetRequiredParam(w, r, 3)
	reaction.Time = &t
	if reaction.Type == 0 {
		reaction.Type = 1
	}
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &reaction)
		if !sv.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Create) {
			result, er3 := h.service.Insert(r.Context(), &reaction)
			sv.AfterCreated(w, r, &reaction, result, er3, h.Status, h.Error, h.Log, h.Resource, h.Action.Create)
		}
	}
}

func (h *reactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var reaction Reaction
	reaction.UserId = sv.GetRequiredParam(w, r, 0)
	reaction.Author = sv.GetRequiredParam(w, r, 2)
	reaction.Id = sv.GetRequiredParam(w, r, 3)

	result, err := h.service.Delete(r.Context(), &reaction)
	sv.HandleDelete(w, r, result, err, h.Error, h.Log, h.Resource, h.Action.Delete)
}

package follow

import (
	"context"
	"net/http"
	"reflect"

	sv "github.com/core-go/core"

)

type FollowHandler interface {
	Check(w http.ResponseWriter, r *http.Request)
	Follow(w http.ResponseWriter, r *http.Request)
	UnFollow(w http.ResponseWriter, r *http.Request)
}

func NewFollowHandler(service FollowService, status sv.StatusConfig, logError func(context.Context, string, ...map[string]interface{}), validate func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error), action *sv.ActionConfig) FollowHandler {
	modelType := reflect.TypeOf(Follower{})
	params := sv.CreateParams(modelType, &status, logError, validate, action)
	return &followHandler{service: service, Params: params}
}

type followHandler struct {
	service FollowService
	*sv.Params
}

func (h *followHandler) Follow(w http.ResponseWriter, r *http.Request) {
	var follower Follower
	follower.Follower = sv.GetRequiredParam(w, r)
	follower.Id = sv.GetRequiredParam(w, r, 1)
	if len(follower.Id) > 0 && len(follower.Follower) > 0 {
		result, err := h.service.Follow(r.Context(), follower.Id, follower.Follower)
		sv.AfterCreated(w, r, &follower, result, err, h.Status, h.Error, h.Log, h.Resource, h.Action.Create)
	}
}

func (h *followHandler) UnFollow(w http.ResponseWriter, r *http.Request) {
	target := sv.GetRequiredParam(w, r)
	id := sv.GetRequiredParam(w, r, 1)
	if len(id) > 0 && len(target) > 0 {
		result, err := h.service.UnFollow(r.Context(), id, target)
		sv.HandleDelete(w, r, result, err, h.Error, h.Log, h.Resource, h.Action.Delete)
	}
}

func (h *followHandler) Check(w http.ResponseWriter, r *http.Request) {
	target := sv.GetRequiredParam(w, r)
	id := sv.GetRequiredParam(w, r, 1)
	if len(id) > 0 && len(target) > 0 {
		result, err := h.service.CheckFollow(r.Context(), id, target)
		sv.RespondModel(w, r, result, err, h.Error, nil)
	}
}
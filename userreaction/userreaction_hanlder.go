package userreaction

import (
	"context"
	"net/http"
	"reflect"

	sv "github.com/core-go/core"
)

type UserReactionHandler interface {
	React(w http.ResponseWriter, r *http.Request)
	Unreact(w http.ResponseWriter, r *http.Request)
	CheckReact(w http.ResponseWriter, r *http.Request)
}

func NewUserReactionHandler(service UserReactionService, status sv.StatusConfig, logError func(context.Context, string, ...map[string]interface{}),
	validate func(context.Context, interface{}) ([]sv.ErrorMessage, error), action sv.ActionConfig) UserReactionHandler {
	params := sv.CreateParams(reflect.TypeOf(UserReaction{}), &status, logError, validate, &action)
	return &userReactionHandler{service: service, Params: params}
}

type userReactionHandler struct {
	*sv.Params
	service UserReactionService
}

func (h *userReactionHandler) CheckReact(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r, 1)
	author := sv.GetRequiredParam(w, r)
	if len(id) > 0 && len(author) > 0 {
		res, err := h.service.CheckReaction(r.Context(), id, author)
		sv.RespondModel(w, r, res, err, h.Error, h.Log)
	}
}

func (h *userReactionHandler) Unreact(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r, 2)
	author := sv.GetRequiredParam(w, r, 1)
	reaction := sv.GetRequiredParam(w, r, 0)
	if len(id) > 0 && len(author) > 0 && len(reaction) > 0 {
		res, err := h.service.Unreact(r.Context(), id, author, reaction)
		sv.RespondModel(w, r, res, err, h.Error, h.Log)
	}
}

func (h *userReactionHandler) React(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r, 2)
	author := sv.GetRequiredParam(w, r, 1)
	reaction := sv.GetRequiredParam(w, r, 0)
	if len(id) > 0 && len(author) > 0 && len(reaction) > 0 {
		res, err := h.service.React(r.Context(), id, author, reaction)
		sv.RespondModel(w, r, res, err, h.Error, h.Log)
	}
}

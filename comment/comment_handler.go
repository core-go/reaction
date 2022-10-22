package comment

import (
	"context"
	"net/http"
	"reflect"
	"time"

	sv "github.com/core-go/core"
	"github.com/google/uuid"
)

type CommentHandler interface {
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewCommentHandler(
	service CommentService,
	status sv.StatusConfig,
	logError func(context.Context, string, ...map[string]interface{}),
	validate func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error),
	action *sv.ActionConfig,
) CommentHandler {
	modelType := reflect.TypeOf(Comment{})
	params := sv.CreateParams(modelType, &status, logError, validate, action)
	return &commentHandler{service: service, Params: params}
}

type commentHandler struct {
	service CommentService
	*sv.Params
}

func (h *commentHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r, 2)
	author := sv.GetRequiredParam(w, r, 1)

	if len(id) > 0 {
		result, err := h.service.Load(r.Context(), id, author)
		sv.RespondModel(w, r, result, err, h.Error, nil)
	}
}
func (h *commentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var comment Comment
	var t = time.Now()
	er1 := sv.Decode(w, r, &comment)
	comment.CommentId = uuid.New().String()
	comment.UserId = sv.GetRequiredParam(w, r)
	comment.Author = sv.GetRequiredParam(w, r, 2)
	comment.Id = sv.GetRequiredParam(w, r, 3)
	comment.Time = &t

	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &comment)
		if !sv.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Create) {
			result, er3 := h.service.Create(r.Context(), &comment)
			sv.AfterCreated(w, r, &comment, result, er3, h.Status, h.Error, h.Log, h.Resource, h.Action.Create)
		}
	}
}

func (h *commentHandler) Update(w http.ResponseWriter, r *http.Request) {
	var comment Comment
	comment.CommentId = sv.GetRequiredParam(w, r)
	comment.UserId = sv.GetRequiredParam(w, r, 1)
	comment.Author = sv.GetRequiredParam(w, r, 3)
	comment.Id = sv.GetRequiredParam(w, r, 4)

	er1 := sv.DecodeAndCheckId(w, r, &comment, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &comment)
		if !sv.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Update) {
			result, er3 := h.service.Update(r.Context(), &comment)
			sv.HandleResult(w, r, &comment, result, er3, h.Status, h.Error, h.Log, h.Resource, h.Action.Update)
		}
	}
}

func (h *commentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var commentId, id, author string
	commentId = sv.GetRequiredParam(w, r)
	author = sv.GetRequiredParam(w, r, 2)
	id = sv.GetRequiredParam(w, r, 3)

	if len(commentId) > 0 {
		result, err := h.service.Delete(r.Context(), commentId, id, author)
		sv.HandleDelete(w, r, result, err, h.Error, h.Log, h.Resource, h.Action.Delete)
	}
}

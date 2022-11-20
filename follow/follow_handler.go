package follow

import (
	"encoding/json"
	"net/http"
	"strings"
)

func NewFollowHandler(service FollowService, targetIndex int, idIndex int) FollowHandler {
	return FollowHandler{service: service, idIndex: idIndex, targetIndex: targetIndex}
}

type FollowHandler struct {
	service     FollowService
	targetIndex int
	idIndex     int
}

func (h *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
	var follower Follower
	follower.Follower = GetRequiredParam(w, r, h.targetIndex)
	follower.Id = GetRequiredParam(w, r, h.idIndex)
	if len(follower.Id) > 0 && len(follower.Follower) > 0 {
		result, err := h.service.Follow(r.Context(), follower.Id, follower.Follower)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
		return
	}

}

func (h *FollowHandler) UnFollow(w http.ResponseWriter, r *http.Request) {
	target := GetRequiredParam(w, r, h.targetIndex)
	id := GetRequiredParam(w, r, h.idIndex)
	if len(id) > 0 && len(target) > 0 {
		result, err := h.service.UnFollow(r.Context(), id, target)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
		return
	}
}

func (h *FollowHandler) Check(w http.ResponseWriter, r *http.Request) {
	target := GetRequiredParam(w, r, h.targetIndex)
	id := GetRequiredParam(w, r, h.idIndex)
	if len(id) > 0 && len(target) > 0 {
		result, err := h.service.CheckFollow(r.Context(), id, target)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
		return
	}
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

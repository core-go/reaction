package rate

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NewRateHandler(
	service RateService,
	authorIndex int,
	idIndex int,
	max int,
) RateHandler {
	return RateHandler{
		service:     service,
		authorIndex: authorIndex,
		idIndex:     idIndex,
		max:max,
	}
}

type RateHandler struct {
	service     RateService
	authorIndex int
	idIndex     int
	max int
}


func (h *RateHandler) Rate(w http.ResponseWriter, r *http.Request) {
	var rate Rate
	var t = time.Now()
	rate.Time = &t
	er1 := Decode(w, r, &rate)
	rate.Author = GetRequiredParam(w, r, h.authorIndex)
	rate.Id = GetRequiredParam(w, r, h.idIndex)

	if er1 == nil {
		errors := Validate(r.Context(), &rate, h.max)
		if len(errors) > 0 {
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(errors)
			return
		}
		result, er3 := h.service.Rate(r.Context(), &rate)
		if er3 != nil {
			http.Error(w, er3.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
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

func Validate(ctx context.Context, rate *Rate, max int) []ErrorMessage {
	errors := []ErrorMessage{}
	if rate.Rate > max {
		errors = append(errors, ErrorMessage{
			Field:   "rate",
			Code:    "max",
			Param: strconv.Itoa(max),
		})
	}
	return errors
}

type ErrorMessage struct {
	Field   string `yaml:"field" mapstructure:"field" json:"field,omitempty" gorm:"column:field" bson:"field,omitempty" dynamodbav:"field,omitempty" firestore:"field,omitempty"`
	Code    string `yaml:"code" mapstructure:"code" json:"code,omitempty" gorm:"column:code" bson:"code,omitempty" dynamodbav:"code,omitempty" firestore:"code,omitempty"`
	Param   string `yaml:"param" mapstructure:"param" json:"param,omitempty" gorm:"column:param" bson:"param,omitempty" dynamodbav:"param,omitempty" firestore:"param,omitempty"`
	Message string `yaml:"message" mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}

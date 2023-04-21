package rates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func NewRatesHandler(
	service RatesService,
	authorIndex int,
	idIndex int,
	max int,
) RatesHandler {
	return RatesHandler{
		service:     service,
		max:         max,
		idIndex:     idIndex,
		authorIndex: authorIndex,
	}
}

type RatesHandler struct {
	service     RatesService
	authorIndex int
	idIndex     int
	max         int
}

func (h *RatesHandler) Rate(w http.ResponseWriter, r *http.Request) {
	var req Request
	er1 := Decode(w, r, &req)
	author := GetRequiredParam(w, r, h.authorIndex) //0
	id := GetRequiredParam(w, r, h.idIndex)         //1
	if er1 == nil {
		errors, er2 := validate(&req, h.max)
		if er2 != nil {
			http.Error(w, er2.Error(), 500)
			return
		}
		if len(errors) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(errors)
			return
		}
		result, er3 := h.service.Rate(r.Context(), id, author, &req)
		fmt.Println(result)
		if er3 != nil {
			fmt.Println("Error : " + er3.Error())
			http.Error(w, er3.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

func validate(rate *Request, max int) ([]ErrorMessage, error) {
	errors := []ErrorMessage{}
	if rate.Rate > float32(max) {
		errors = append(errors, ErrorMessage{Field: "rate", Code: "max", Param: strconv.Itoa(max)})
	}
	if len(rate.Rates) == 0 {
		errors = append(errors, ErrorMessage{Field: "rate", Code: "required"})
	}
	for i := 0; i < len(rate.Rates); i++ {
		if rate.Rates[i] > float32(max) {
			errors = append(errors, ErrorMessage{Field: fmt.Sprintf("rate$1", i), Code: "max", Param: strconv.Itoa(max)})
		}
	}
	return errors, nil
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

type ErrorMessage struct {
	Field   string `yaml:"field" mapstructure:"field" json:"field,omitempty" gorm:"column:field" bson:"field,omitempty" dynamodbav:"field,omitempty" firestore:"field,omitempty"`
	Code    string `yaml:"code" mapstructure:"code" json:"code,omitempty" gorm:"column:code" bson:"code,omitempty" dynamodbav:"code,omitempty" firestore:"code,omitempty"`
	Param   string `yaml:"param" mapstructure:"param" json:"param,omitempty" gorm:"column:param" bson:"param,omitempty" dynamodbav:"param,omitempty" firestore:"param,omitempty"`
	Message string `yaml:"message" mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}

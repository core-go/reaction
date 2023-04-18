package rate

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Request struct {
	Rate      int    `json:"rate,omitempty" gorm:"column:rate" bson:"rate,omitempty" dynamodbav:"rate,omitempty" firestore:"rate,omitempty" validate:"required,max=10"`
	Review    string `json:"review,omitempty" gorm:"column:review" bson:"review,omitempty" dynamodbav:"review,omitempty" firestore:"review,omitempty"`
	Anonymous bool   `json:"anonymous,omitempty" gorm:"column:anonymous" bson:"anonymous,omitempty" dynamodbav:"anonymous,omitempty" firestore:"anonymous,omitempty"`
}

type Rate struct {
	Id          string      `json:"id,omitempty" gorm:"column:id;primary_key" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty" validate:"required,max=255"`
	Author      string      `json:"author,omitempty" gorm:"column:author;primary_key" bson:"author,omitempty" dynamodbav:"author,omitempty" firestore:"author,omitempty" validate:"required,max=255"`
	Rate        int         `json:"rate,omitempty" gorm:"column:rate" bson:"rate,omitempty" dynamodbav:"rate,omitempty" firestore:"rate,omitempty" validate:"required,max=10"`
	Review      string      `json:"review,omitempty" gorm:"column:review" bson:"review,omitempty" dynamodbav:"review,omitempty" firestore:"review,omitempty"`
	Time        *time.Time  `json:"time,omitempty" gorm:"column:time" bson:"time,omitempty" dynamodbav:"time,omitempty" firestore:"time,omitempty"`
	UsefulCount int         `json:"usefulCount,omitempty" gorm:"column:usefulCount" bson:"usefulCount,omitempty" dynamodbav:"usefulCount,omitempty" firestore:"usefulCount,omitempty"`
	ReplyCount  int         `json:"replyCount,omitempty" gorm:"column:replyCount" bson:"replyCount,omitempty" dynamodbav:"replyCount,omitempty" firestore:"replyCount,omitempty"`
	Histories   []Histories `json:"histories,omitempty" gorm:"column:histories" bson:"histories,omitempty" dynamodbav:"histories,omitempty" firestore:"histories,omitempty"`
	Anonymous   bool        `json:"anonymous,omitempty" gorm:"column:anonymous" bson:"anonymous,omitempty" dynamodbav:"anonymous,omitempty" firestore:"anonymous,omitempty"`
}

type RateInfo struct {
	Id     string  `gorm:"column:id;primary_key" validate:"required,max=255"`
	Rate   float32 `gorm:"column:rate;"`
	Rate1  int     `gorm:"column:rate1;"`
	Rate2  int     `gorm:"column:rate2;"`
	Rate3  int     `gorm:"column:rate3;"`
	Rate4  int     `gorm:"column:rate4;"`
	Rate5  int     `gorm:"column:rate5;"`
	Rate6  int     `gorm:"column:rate6;"`
	Rate7  int     `gorm:"column:rate7;"`
	Rate8  int     `gorm:"column:rate8;"`
	Rate9  int     `gorm:"column:rate9;"`
	Rate10 int     `gorm:"column:rate10;"`
	Count  int     `gorm:"column:count;"`
	Score  int     `gorm:"column:score;"`
}

type Histories struct {
	Time   *time.Time `json:"time,omitempty" gorm:"column:time" bson:"time,omitempty" dynamodbav:"time,omitempty" firestore:"time,omitempty" validate:"required"`
	Rate   int        `json:"rate,omitempty" gorm:"column:rate" bson:"rate,omitempty" dynamodbav:"rate,omitempty" firestore:"rate,omitempty" validate:"required"`
	Review string     `json:"review,omitempty" gorm:"column:review" bson:"review,omitempty" dynamodbav:"review,omitempty" firestore:"review,omitempty"`
}

func (c Histories) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Histories) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

type ErrorMessage struct {
	Field   string `yaml:"field" mapstructure:"field" json:"field,omitempty" gorm:"column:field" bson:"field,omitempty" dynamodbav:"field,omitempty" firestore:"field,omitempty"`
	Code    string `yaml:"code" mapstructure:"code" json:"code,omitempty" gorm:"column:code" bson:"code,omitempty" dynamodbav:"code,omitempty" firestore:"code,omitempty"`
	Param   string `yaml:"param" mapstructure:"param" json:"param,omitempty" gorm:"column:param" bson:"param,omitempty" dynamodbav:"param,omitempty" firestore:"param,omitempty"`
	Message string `yaml:"message" mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}

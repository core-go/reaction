package searchresponse

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Response struct {
	Id           string      `json:"id,omitempty" gorm:"column:id;primary_key" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty" validate:"required,max=255"`
	Author       string      `json:"author,omitempty" gorm:"column:author;primary_key" bson:"author,omitempty" dynamodbav:"author,omitempty" firestore:"author,omitempty" validate:"required,max=255"`
	Description  string      `json:"description,omitempty" gorm:"column:description" bson:"description,omitempty" dynamodbav:"description,omitempty" firestore:"description,omitempty"`
	Time         *time.Time  `json:"time,omitempty" gorm:"column:time" bson:"time,omitempty" dynamodbav:"time,omitempty" firestore:"time,omitempty"`
	UsefulCount  int         `json:"usefulCount,omitempty" gorm:"column:usefulCount" bson:"usefulCount,omitempty" dynamodbav:"usefulCount,omitempty" firestore:"usefulCount,omitempty"`
	CommentCount int         `json:"commentCount,omitempty" gorm:"column:commentCount" bson:"commentCount,omitempty" dynamodbav:"commentCount,omitempty" firestore:"commentCount,omitempty"`
	Histories    []Histories `json:"histories,omitempty" gorm:"column:histories" bson:"histories,omitempty" dynamodbav:"histories,omitempty" firestore:"histories,omitempty"`
}

type ResponseInfo struct {
	Id    string `gorm:"column:id;primary_key" validate:"required,max=255"`
	Count int    `gorm:"column:count;"`
}

type Histories struct {
	Time        *time.Time `json:"time,omitempty" gorm:"column:time" bson:"time,omitempty" dynamodbav:"time,omitempty" firestore:"time,omitempty" validate:"required"`
	Description string     `json:"description,omitempty" gorm:"column:description" bson:"description,omitempty" dynamodbav:"description,omitempty" firestore:"description,omitempty"`
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

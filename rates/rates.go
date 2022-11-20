package rates

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Rates struct {
	Id          string      `json:"id" gorm:"id"`
	Author      string      `json:"author" gorm:"author"`
	Rate        float32     `json:"rate" gorm:"rate"`
	Rates       []float32   `json:"rates" gorm:"rates"`
	Time        *time.Time  `json:"time" gorm:"time"`
	Review      string      `json:"review" gorm:"review"`
	UsefulCount int         `json:"usefulcount" gorm:"usefulcount"`
	ReplyCount  int         `json:"replycount" gorm:"replycount"`
	Histories   []Histories `json:"histories" gorm:"histories"`
}

type Histories struct {
	Time   *time.Time `json:"time,omitempty" gorm:"column:time" bson:"time,omitempty" dynamodbav:"time,omitempty" firestore:"time,omitempty" validate:"required"`
	Rate   float32    `json:"rate,omitempty" gorm:"column:rate" bson:"rate,omitempty" dynamodbav:"rate,omitempty" firestore:"rate,omitempty" validate:"required"`
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

package save

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Save struct {
	Id   string   `json:"id,omitempty" gorm:"column:id;primary_key" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty" validate:"required,max=255"`
	Item []string `json:"item,omitempty" gorm:"column:item" bson:"item,omitempty" dynamodbav:"item,omitempty" firestore:"item,omitempty" validate:"required"`
}

func (c Save) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Save) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

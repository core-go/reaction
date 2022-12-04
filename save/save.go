package save

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Items struct {
	Id    string   `json:"id,omitempty" gorm:"column:id;primary_key" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty" validate:"required,max=255" match:"equal"`
	Items []string `json:"items,omitempty" gorm:"column:items" bson:"items,omitempty" dynamodbav:"items,omitempty" firestore:"items,omitempty" validate:"required"`
}

func (c Items) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Items) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

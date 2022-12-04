package comment

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Comment struct {
	CommentId string      `json:"commentId,omitempty" gorm:"column:commentId;primary_key" bson:"_commentId,omitempty" dynamodbav:"commentId,omitempty" firestore:"commentId,omitempty" validate:"required,max=40"`
	Id        string      `json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty" validate:"required,max=255"`
	Author    string      `json:"author,omitempty" gorm:"column:author" bson:"author,omitempty" dynamodbav:"author,omitempty" firestore:"author,omitempty" validate:"required,max=255"`
	UserId    string      `json:"userId,omitempty" gorm:"column:userId" bson:"userId,omitempty" dynamodbav:"userId,omitempty" firestore:"userId,omitempty" validate:"required,max=255"`
	Comment   string      `json:"comment,omitempty" gorm:"column:comment" bson:"comment,omitempty" dynamodbav:"comment,omitempty" firestore:"comment,omitempty"`
	Time      *time.Time  `json:"time,omitempty" gorm:"column:time" bson:"time,omitempty" dynamodbav:"time,omitempty" firestore:"time,omitempty"`
	UpdatedAt *time.Time  `json:"updateAt,omitempty" gorm:"column:updateAt" bson:"updateAt,omitempty" dynamodbav:"updateAt,omitempty" firestore:"updateAt,omitempty"`
	Histories []Histories `json:"histories,omitempty" gorm:"column:histories" bson:"histories,omitempty" dynamodbav:"histories,omitempty" firestore:"histories,omitempty"`
	UserURL   *string     `json:"userURL,omitempty" gorm:"column:imageurl"`
	Username  *string     `json:"authorName,omitempty" gorm:"column:username"`
}

type Histories struct {
	Time    *time.Time `json:"time,omitempty" gorm:"column:time" bson:"time,omitempty" dynamodbav:"time,omitempty" firestore:"time,omitempty"`
	Comment string     `json:"comment,omitempty" gorm:"column:comment" bson:"comment,omitempty" dynamodbav:"comment,omitempty" firestore:"comment,omitempty"`
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

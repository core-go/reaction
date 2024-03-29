package comment

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Request struct {
	Comment   string `json:"comment,omitempty" gorm:"column:comment" bson:"comment,omitempty" dynamodbav:"comment,omitempty" firestore:"comment,omitempty"`
	Anonymous bool   `json:"anonymous,omitempty" gorm:"column:anonymous" bson:"anonymous,omitempty" dynamodbav:"anonymous,omitempty" firestore:"anonymous,omitempty"`
}

type Response struct {
	CommentId  string      `json:"commentId,omitempty" gorm:"column:commentId;primary_key" bson:"_commentId,omitempty" dynamodbav:"commentId,omitempty" firestore:"commentId,omitempty" validate:"required,max=40"`
	Id         string      `json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty" validate:"required,max=255"`
	Author     string      `json:"author,omitempty" gorm:"column:author" bson:"author,omitempty" dynamodbav:"author,omitempty" firestore:"author,omitempty" validate:"required,max=255"`
	UserId     string      `json:"userId,omitempty" gorm:"column:userId" bson:"userId,omitempty" dynamodbav:"userId,omitempty" firestore:"userId,omitempty" validate:"required,max=255"`
	Comment    string      `json:"comment,omitempty" gorm:"column:comment" bson:"comment,omitempty" dynamodbav:"comment,omitempty" firestore:"comment,omitempty"`
	Anonymous  bool        `json:"anonymous,omitempty" gorm:"column:anonymous" bson:"anonymous,omitempty" dynamodbav:"anonymous,omitempty" firestore:"anonymous,omitempty"`
	Time       *time.Time  `json:"time,omitempty" gorm:"column:time" bson:"time,omitempty" dynamodbav:"time,omitempty" firestore:"time,omitempty"`
	UpdatedAt  *time.Time  `json:"updateAt,omitempty" gorm:"column:updateAt" bson:"updateAt,omitempty" dynamodbav:"updateAt,omitempty" firestore:"updateAt,omitempty"`
	Histories  []Histories `json:"histories,omitempty" gorm:"column:histories" bson:"histories,omitempty" dynamodbav:"histories,omitempty" firestore:"histories,omitempty"`
	AuthorURL  *string     `json:"authorURL,omitempty" gorm:"column:imageurl"`
	AuthorName *string     `json:"authorName,omitempty" gorm:"column:username"`
}

func toResponse(c Comment) Response {
	return Response{
		CommentId: c.CommentId,
		Id:        c.Id,
		Author:    c.Author,
		UserId:    c.UserId,
		Comment:   c.Comment,
		Anonymous: c.Anonymous,
		Time:      c.Time,
		UpdatedAt: c.UpdatedAt,
		Histories: c.Histories,
	}
}

type Comment struct {
	CommentId string      `json:"commentId,omitempty" gorm:"column:commentId;primary_key" bson:"_commentId,omitempty" dynamodbav:"commentId,omitempty" firestore:"commentId,omitempty" validate:"required,max=40"`
	Id        string      `json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty" validate:"required,max=255"`
	Author    string      `json:"author,omitempty" gorm:"column:author" bson:"author,omitempty" dynamodbav:"author,omitempty" firestore:"author,omitempty" validate:"required,max=255"`
	UserId    string      `json:"userId,omitempty" gorm:"column:userId" bson:"userId,omitempty" dynamodbav:"userId,omitempty" firestore:"userId,omitempty" validate:"required,max=255"`
	Comment   string      `json:"comment,omitempty" gorm:"column:comment" bson:"comment,omitempty" dynamodbav:"comment,omitempty" firestore:"comment,omitempty"`
	Anonymous bool        `json:"anonymous,omitempty" gorm:"column:anonymous" bson:"anonymous,omitempty" dynamodbav:"anonymous,omitempty" firestore:"anonymous,omitempty"`
	Time      *time.Time  `json:"time,omitempty" gorm:"column:time" bson:"time,omitempty" dynamodbav:"time,omitempty" firestore:"time,omitempty"`
	UpdatedAt *time.Time  `json:"updateAt,omitempty" gorm:"column:updateAt" bson:"updateAt,omitempty" dynamodbav:"updateAt,omitempty" firestore:"updateAt,omitempty"`
	Histories []Histories `json:"histories,omitempty" gorm:"column:histories" bson:"histories,omitempty" dynamodbav:"histories,omitempty" firestore:"histories,omitempty"`
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

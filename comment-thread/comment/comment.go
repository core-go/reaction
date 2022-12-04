package comment

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Comment struct {
	CommentId       string     `json:"commentId" gorm:"column:commentid" bson:"commentId" firestore:"commentId"`
	Id              string     `json:"id" gorm:"column:id" bson:"id" firestore:"id"`
	Author          string     `json:"author" gorm:"column:author" bson:"author" firestore:"author"`
	UserId          string     `json:"userId" gorm:"column:userid" bson:"userId" firestore:"userId"`
	Comment         string     `json:"comment" gorm:"column:comment" bson:"comment" firestore:"comment"`
	Parent          *string    `json:"parent" gorm:"column:parent" bson:"parent" firestore:"parent"`
	Time            time.Time  `json:"time" gorm:"column:time" bson:"time" firestore:"time"`
	CommentThreadId string     `json:"commentThreadId" gorm:"commentthreadid" bson:"commentThreadId" firestore:"commentThreadId"`
	UpdatedAt       *time.Time `json:"updatedAt" gorm:"column:updatedat" bson:"updatedAt" firestore:"updatedAt"`
	Histories       []History  `json:"histories" gorm:"column:histories" bson:"histories" firestore:"histories"`
	ReplyCount      *int       `json:"replyCount" gorm:"column:replycount" bson:"replyCount" firestore:"replyCount"`
	UsefulCount     *int       `json:"usefulCount" gorm:"column:usefulcount" bson:"usefulCount" firestore:"usefulCount"`
	Username        *string    `json:"authorName" gorm:"column:username" bson:"authorName" firestore:"authorName"`
	Avatar          *string    `json:"userURL" gorm:"column:avatar" bson:"userURL" firestore:"userURL"`
	Disable         *bool      `json:"disable" gorm:"column:-" bson:"-" firestore:"-"`
}

func (c History) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *History) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

type History struct {
	Comment string    `json:"comment" gorm:"column:comment"`
	Time    time.Time `json:"time" gorm:"column:time"`
}

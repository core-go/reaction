package comment

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Comment struct {
	CommentId       string     `json:"commentId" gorm:"column:commentid"`
	Id              string     `json:"id" gorm:"column:id"`
	Author          string     `json:"author" gorm:"column:author"`
	Comment         string     `json:"comment" gorm:"column:comment"`
	Time            time.Time  `json:"time" gorm:"column:time"`
	CommentThreadId string     `json:"commentThreadId" gorm:"commentthreadid"`
	UpdatedAt       *time.Time `json:"updatedat" gorm:"column:updatedat"`
	Histories       []History  `json:"histories" gorm:"column:histories"`
	ReplyCount      *int       `json:"replyCount" gorm:"column:replycount"`
	UsefulCount     *int       `json:"usefulCount" gorm:"column:usefulcount"`
	Disable         *bool      `json:"disable" gorm:"column:disable"`
}

type Request struct {
	Comment string `json:"comment" gorm:"column:comment"`
}

type Response struct {
	CommentId       string     `json:"commentId" gorm:"column:commentid"`
	Id              string     `json:"id" gorm:"column:id"`
	Author          string     `json:"author" gorm:"column:author"`
	Comment         string     `json:"comment" gorm:"column:comment"`
	Time            time.Time  `json:"time" gorm:"column:time"`
	CommentThreadId string     `json:"commentThreadId" gorm:"commentthreadid"`
	UpdatedAt       *time.Time `json:"updatedat" gorm:"column:updatedat"`
	Histories       []History  `json:"histories" gorm:"column:histories"`
	ReplyCount      *int       `json:"replyCount" gorm:"column:replycount"`
	UsefulCount     *int       `json:"usefulCount" gorm:"column:usefulcount"`
	AuthorName      *string    `json:"authorName" gorm:"column:username"`
	AuthorURL       *string    `json:"authorURL" gorm:"column:imageurl"`
	Disable         *bool      `json:"disable" gorm:"column:disable"`
}

func toResponse(c Comment) Response {
	return Response{
		CommentId:       c.CommentId,
		Id:              c.Id,
		Author:          c.Author,
		Comment:         c.Comment,
		Time:            c.Time,
		CommentThreadId: c.CommentThreadId,
		UpdatedAt:       c.UpdatedAt,
		Histories:       c.Histories,
		ReplyCount:      c.ReplyCount,
		UsefulCount:     c.UsefulCount,
		Disable:         c.Disable,
	}
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

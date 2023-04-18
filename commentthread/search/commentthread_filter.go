package search

import (
	"time"

	"github.com/core-go/search"
)

type CommentThreadFilter struct {
	*search.Filter
	CommentId string     `json:"commentId" gorm:"column:commentid;primary_key" bson:"commentId" firestore:"commentId" match:"equal"`
	Id        string     `json:"id" gorm:"column:id" bson:"id" firestore:"id" match:"equal"`
	UserId    string     `json:"userId" gorm:"column:userId" bson:"userId" firestore:"userId"`
	Author    string     `json:"author" gorm:"column:author" bson:"author" firestore:"author" match:"equal"`
	Comment   string     `json:"comment" gorm:"column:comment" bson:"comment" firestore:"comment"`
	Time      time.Time  `json:"time" gorm:"column:time" bson:"time" firestore:"time"`
	UpdatedAt *time.Time `json:"updateAt" gorm:"column:updatedat" bson:"updatedat" firestore:"updatedat"`
}

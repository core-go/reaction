package searchcomment

import (
	"github.com/core-go/search"
)

type CommentFilter struct {
	*search.Filter
	CommentId string            `mapstructure:"commentId" json:"commentId" gorm:"column:commentId;primary_key" bson:"_commentId" dynamodbav:"commentId" firestore:"commentId" match:"equal" validate:"max=40"`
	Id        string            `mapstructure:"id" json:"id" gorm:"column:id" bson:"id" dynamodbav:"id" firestore:"id" match:"equal" validate:"max=255"`
	Author    string            `mapstructure:"author" json:"author" gorm:"column:author" bson:"author" dynamodbav:"author" firestore:"author" match:"equal" validate:"max=255"`
	UserId    string            `mapstructure:"userId" json:"userId" gorm:"column:userId" bson:"userId" dynamodbav:"userId" firestore:"userId" match:"equal" validate:"max=255"`
	Comment   string            `mapstructure:"comment" json:"comment" gorm:"column:comment" bson:"comment" dynamodbav:"comment" firestore:"comment"`
	Time      *search.TimeRange `mapstructure:"time" json:"time" gorm:"column:time" bson:"time" dynamodbav:"time" firestore:"time"`
}

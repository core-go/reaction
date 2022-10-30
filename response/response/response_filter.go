package response

import (
	"github.com/core-go/search"
)

type ResponseFilter struct {
	*search.Filter
	Id           string            `mapstructure:"id" json:"id,omitempty" gorm:"column:id;primary_key" bson:"id" dynamodbav:"id" firestore:"id" match:"equal" validate:"max=255"`
	Author       string            `mapstructure:"author" json:"author,omitempty" gorm:"column:author;primary_key" bson:"author" dynamodbav:"author" firestore:"author" match:"equal" validate:"max=255"`
	Desciption   string            `mapstructure:"description" json:"description,omitempty" gorm:"column:description" bson:"descri" dynamodbav:"descri" firestore:"descri"`
	Time         *search.TimeRange `mapstructure:"time" json:"time" gorm:"column:time" bson:"time" dynamodbav:"time" firestore:"time"`
	UsefulCount  string            `mapstructure:"usefulCount" json:"usefulCount,omitempty" gorm:"column:usefulCount" bson:"usefulCount" dynamodbav:"usefulCount" firestore:"usefulCount"`
	CommentCount string            `mapstructure:"commentCount" json:"commentCount,omitempty" gorm:"column:commentCount" bson:"commentCount" dynamodbav:"commentCount" firestore:"commentCount"`
}

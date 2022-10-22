package reaction

import (
	"time"
)

type Reaction struct {
	Id     string     `json:"id,omitempty" gorm:"column:id;primary_key" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty" validate:"required,max=255"`
	Author string     `json:"author,omitempty" gorm:"column:author;primary_key" bson:"author,omitempty" dynamodbav:"author,omitempty" firestore:"author,omitempty" validate:"required,max=255"`
	UserId string     `json:"userId,omitempty" gorm:"column:userId;primary_key" bson:"userId,omitempty" dynamodbav:"userId,omitempty" firestore:"userId,omitempty" validate:"required,max=255"`
	Time   *time.Time `json:"time,omitempty" gorm:"column:time" bson:"time,omitempty" dynamodbav:"time,omitempty" firestore:"time,omitempty"`
	Type   int8       `json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty" validate:"max=10"`
}

package comment

type CommentFilter struct {
	CommentId string `json:"commentId" gorm:"column:commentid" bson:"commentId" firestore:"commentId" match:"equal"`
	Id        string `json:"id" gorm:"column:id" bson:"id" firestore:"id" match:"equal"`
	Author    string `json:"author" gorm:"column:author" bson:"author" firestore:"author" match:"equal"`
	UserId    string `json:"userId" gorm:"column:userid" bson:"userId" firestore:"userId" match:"equal"`
	Comment   string `json:"comment" gorm:"column:comment" bson:"comment" firestore:"comment"`
}

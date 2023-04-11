package comment

type CommentFilter struct {
	CommentId string `json:"commentId" gorm:"column:commentid" match:"equal"`
	Id        string `json:"id" gorm:"column:id" match:"equal"`
	Author    string `json:"author" gorm:"column:author" match:"equal"`
	UserId    string `json:"userId" gorm:"column:userid" match:"equal"`
	Comment   string `json:"comment" gorm:"column:comment"`
}

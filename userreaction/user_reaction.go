package userreaction

type UserReaction struct {
	Id       string `json:"id,omitempty" gorm:"column:id,primary_key"`
	Author   string `json:"author,omitempty" gorm:"column:author"`
	Reaction int64  `json:"reaction,omitempty" gorm:"column:reaction"`
}

package search

type Info struct {
	Id   string `json:"id,omitempty" gorm:"column:id;primary_key"`
	Url  string `json:"url,omitempty" gorm:"column:url"`
	Name string `json:"name,omitempty" gorm:"column:name"`
}

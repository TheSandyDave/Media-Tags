package domain

type Media struct {
	BaseObject
	Name    string
	Tags    []*Tag `gorm:"many2many:media_tags;constaint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FileUrl string
}

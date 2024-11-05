package domain

type Tag struct {
	BaseObject
	Name string `gorm:"unique"`
}

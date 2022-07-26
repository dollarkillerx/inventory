package models

type UserCenter struct {
	BasicModel
	Storehouse string `gorm:"type:varchar(300)" json:"storehouse"`
	Account    string `gorm:"type:varchar(300);index" json:"account"`
	Password   string `gorm:"type:varchar(600)" json:"password"`
}

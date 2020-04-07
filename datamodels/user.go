package datamodels

type User struct {
	Id int64 `json:"id" form:"id" sql:"id"`
	Nickname string `json:"nickname" form:"nickname" sql:"nickname"`
    UserName string `json:"userName" form:"userName" sql:"userName"`
	HashPassword string `json:"-" form:"passWord" sql:"passWord"`

}

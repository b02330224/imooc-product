package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"imooc-product/datamodels"
	"imooc-product/services"
	"strconv"
)

type UserController struct {
	Ctx iris.Context
	Service services.IUserService
	Session *sessions.Session
}

func (c *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name:"user/register.html",
	}
}

func (c *UserController) PostRegister() {
	var (
		nickName = c.Ctx.FormValue("nickname")
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("password")
	)

	//ozzo-validation 验证表单

	user := &datamodels.User{
		Nickname:     nickName,
		UserName:     userName,
		HashPassword: password,
	}

	_, err := c.Service.AddUser(user)
	if err != nil {
		c.Ctx.Redirect("/user/error")
		return
	}

	c.Ctx.Redirect("/user/login")
	return
}

func (c *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name : "user/login.html",
	}
}

func (c *UserController) PostLogin() mvc.Response {
	var (
		userName = c.Ctx.FormValue("userName" )
		password = c.Ctx.FormValue("password")
	)

	user, isOk := c.Service.IsPwdSuccess(userName, password)
	if !isOk {
		return mvc.Response{
			Path:"/user/login",
		}
	}

	c.Ctx.SetCookieKV("userId", strconv.FormatInt(user.Id, 10))
	c.Session.Set("userId", strconv.FormatInt(user.Id, 10))
	return mvc.Response{
		Path:"/product",
	}
}
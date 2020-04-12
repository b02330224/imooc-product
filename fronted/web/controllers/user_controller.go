package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"imooc-product/datamodels"
	"imooc-product/encrypt"
	"imooc-product/services"
	"strconv"
)

type UserController struct {
	Ctx iris.Context
	Service services.IUserService
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

	uidByte := []byte(strconv.FormatInt(user.Id, 10))
	uidString, err := encrypt.EnPwdCode(uidByte)
	if err != nil {
		fmt.Println(err)
	}

	c.Ctx.SetCookieKV("uid", strconv.FormatInt(user.Id, 10), iris.CookieHTTPOnly(false))
	c.Ctx.SetCookieKV("sign", uidString, iris.CookieHTTPOnly(false))

	return mvc.Response{
		Path:"/product",
	}
}
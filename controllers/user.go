package controllers

import (
	"Article/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserController struct {
	beego.Controller
}

// 展示注册页面
func (this *UserController) ShowRegister() {
	this.TplName = "register.html"
}

// 注册用户
func (this *UserController) HandleRegister() {
	// 获取数据
	userName := this.GetString("userName")
	password := this.GetString("password")
	// 校验数据
	if "" == userName {
		this.Data["userErr"] = "用户名不能为空"
		this.TplName = "register.html"
		return
	}

	if "" == password {
		this.Data["passErr"] = "密码不能为空"
		this.TplName = "register.html"
		return
	}

	// 处理数据
	o := orm.NewOrm()

	var user models.User
	user.UserName = userName
	user.PassWord = password

	_, err := o.Insert(&user)
	if err != nil {
		beego.Error("新建用户失败 : ", err)
		this.Data["errMsg"] = "注册失败"
		this.TplName = "register.html"
		return
	}
	// 返回视图
	this.Redirect("/login", 302)
}

// 展示登陆页面
func (this *UserController) ShowLogin() {
	userName := this.Ctx.GetCookie("userName")

	if "" == userName {
		this.Data["userName"] = userName
		this.Data["checked"] = ""
	}else {
		this.Data["userName"] = userName
		this.Data["checked"] = "checked"
	}

	this.TplName = "login.html"
}

// 用户登陆
func (this *UserController) HandleLogin() {
	// 获取数据
	userName := this.GetString("userName")
	password := this.GetString("password")
	remember := this.GetString("remember")
	// 校验数据
	if "" == userName {
		this.Data["userErr"] = "用户名不能为空"
		this.TplName = "login.html"
		return
	}

	if "" == password {
		this.Data["passErr"] = "密码不能为空"
		this.TplName = "login.html"
		return
	}

	// 处理数据
	o := orm.NewOrm()

	var user models.User
	user.UserName = userName

	err := o.Read(&user, "UserName")
	if err != nil {
		this.Data["errMsg"] = "用户名不存在"
		this.TplName = "login.html"
		return
	}

	if user.PassWord != password {
		this.Data["errMsg"] = "密码输入错误"
		this.TplName = "login.html"
		return
	}

	if remember == "on" {
		this.Ctx.SetCookie("userName", userName, 3600*24)
	}else {
		this.Ctx.SetCookie("userName", userName, -1)
	}
	this.SetSession("userName", userName)

	// 返回视图
	//this.TplName = "index.html"
	this.Redirect("/showArticleList", 302)
}

func (this *UserController) Logout() {
	this.DelSession("userName")
	this.Redirect("/login", 302)
}

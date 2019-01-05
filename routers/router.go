package routers

import (
	"Article/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/article/*", beego.BeforeExec, Filter)

	beego.Router("/", &controllers.MainController{})

	// 注册
	beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleRegister")

	// 登陆
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")

	// 显示文章列表页
	beego.Router("/showArticleList", &controllers.ArticleController{}, "get:ShowArticleList")

	// 展示添加文章
	beego.Router("/addArticle", &controllers.ArticleController{}, "get:ShowAddArticle;post:HandleAddArticle")

	// 查看文章行情
	beego.Router("/showArticleContent", &controllers.ArticleController{}, "get:ShowArticleContent")

	// 更新文章
	beego.Router("/updateArticle", &controllers.ArticleController{}, "get:ShowUpdateArticle;post:HandleUpdateArticle")

	// 删除文章
	beego.Router("/delArticle", &controllers.ArticleController{}, "get:DelArticle")

	// 添加文章分类
	beego.Router("/addArticleType", &controllers.ArticleController{}, "get:ShowAddArticleType;post:HandleAddArticleType")

	// 退出登陆
	beego.Router("/logout", &controllers.UserController{}, "get:Logout")

	// 删除类型
	beego.Router("/deleteType", &controllers.ArticleController{}, "get:DelType")
}

var Filter = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
	}
}
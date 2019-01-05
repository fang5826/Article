package controllers

import (
	"Article/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"path"
	"strconv"
	"time"
)

type ArticleController struct {
	beego.Controller
}

// 展示文章列表
func (this *ArticleController) ShowArticleList() {
	typeName := this.GetString("select")
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}

	o := orm.NewOrm()
	var articles []models.Article
	qs := o.QueryTable("Article")
	// 实现分页功能
	// 1. 获取记录数
	var count int64
	var err error
	if typeName == "" {
		count, err = qs.Count()
	} else {
		count, err = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()
	}

	if err != nil {
		beego.Error("获取记录数失败")
		this.Redirect("/showArticleList", 302)
		return
	}

	// 2. 定义每一页显示条数
	pageSize := 2

	// 3. 计算总页数
	pageCount := math.Ceil(float64(count) / float64(pageSize))
	//beego.Info("总条数：", count, "每一页显示条数：", pageSize, "总页数：", pageCount)

	// 4. 获取传入的页码
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}

	// 5.从数据库获取部分数据
	start := pageSize * (pageIndex - 1)
	// 点击下拉框，显示相应的新闻
	if typeName == "" {
		qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)
	} else {
		qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles)
	}

	// 查询所有的类型
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)

	// 返回视图
	this.Data["count"] = count
	this.Data["select"] = typeName
	this.Data["userName"] = userName
	this.Data["articles"] = articles
	this.Data["pageIndex"] = pageIndex
	this.Data["pageCount"] = pageCount
	this.Data["articleTypes"] = articleTypes


	this.TplName = "index.html"
}

// 展示添加文章页面
func (this *ArticleController) ShowAddArticle() {
	// 查询类型
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	_, err := o.QueryTable("ArticleType").All(&articleTypes)
	if err != nil {
		beego.Error("查询失败")
		this.Redirect("/showArticleList", 302)
	}

	this.Data["userName"] = this.GetSession("userName")
	this.Data["articleTypes"] = articleTypes
	this.TplName = "add.html"
}

// 处理添加文章请求
func (this *ArticleController) HandleAddArticle() {
	// 获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	typeName := this.GetString("select")

	// 校验数据
	if "" == articleName {
		this.Data["errmsg"] = "文章名不能为空"
		this.Redirect("/addArticle", 302)
		return
	}

	if "" == content {
		this.Data["errmsg"] = "文章内容不能为空"
		this.Redirect("/addArticle", 302)
		return
	}

	if "" == typeName {
		beego.Error("获取类型失败")
		this.Redirect("/addArticle", 302)
		return
	}

	// 图片上传
	file, head, err := this.GetFile("uploadname")
	if err != nil {
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = "add.html"
		return
	}
	defer file.Close()

	// 检查文件大小
	if head.Size > 500000 {
		this.Data["errmsg"] = "文件太大"
		this.TplName = "add.html"
		return
	}

	// 检查文件类型
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		this.Data["errmsg"] = "文件类型错误， 请重新上传"
		this.TplName = "add.html"
		return
	}

	// 防止重名
	fileName := strconv.Itoa(int(time.Now().UnixNano()))

	//保存图片
	err = this.SaveToFile("uploadname", "./static/img/"+fileName+ext)
	if err != nil {
		beego.Error("文件保存失败", err)
		return
	}

	// 处理数据
	o := orm.NewOrm()

	var article models.Article
	article.Title = articleName
	article.Content = content
	article.Img = "/static/img/" + fileName + ext

	// 获取文章类型，根据文章类型名称获取类型对象， 将该对象赋值给文章的文章类型字段
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType, "TypeName")

	article.ArticleType = &articleType

	o.Insert(&article)

	// 返回视图
	this.Redirect("/showArticleList", 302)
}

// 展示文章内容
func (this *ArticleController) ShowArticleContent() {
	// 获取数据
	articleId, err := this.GetInt("articleId")
	// 校验数据
	if err != nil {
		beego.Error("获取文章ID失败")
		//this.TplName = "index.html"
		this.Redirect("/showArticleList", 302)
		return
	}

	// 数据处理
	o := orm.NewOrm()
	var article models.Article

	article.Id = articleId

	err = o.Read(&article)
	if err != nil {
		beego.Error("获取数据失败")
		// 这里不能直接渲染，因为该页面需要数据，但是没有传入
		//this.TplName = "index.html"
		this.Redirect("/showArticleList", 302)
		return
	}

	// 更新浏览次数
	article.Count += 1
	o.Update(&article)


	// 浏览人
	m2m := o.QueryM2M(&article, "Users")
	var user models.User
	userName := this.GetSession("userName")
	user.UserName = userName.(string)
	err = o.Read(&user, "UserName")
	if err != nil {
		beego.Error("获取失败")
		this.Redirect("/showArticleList", 302)
		return
	}

	_, err = m2m.Add(user)
	if err != nil {
		beego.Error("添加失败")
		this.Redirect("/showArticleList", 302)
		return
	}


	// 获取浏览记录
	var users []models.User
	qs := o.QueryTable("User")
	_, err = qs.Filter("Articles__Article__Id", articleId).Distinct().All(&users)
	if err != nil {
		beego.Error("获取浏览的用户失败")
		this.Redirect("/showArticleList", 302)
		return
	}



	// 返回视图
	this.Data["users"] = users
	this.Data["article"] = article
	this.Data["userName"] = this.GetSession("userName")
	this.TplName = "content.html"
}

// 展示更新文章页面
func (this *ArticleController) ShowUpdateArticle() {
	// 获取数据
	articleId, err := this.GetInt("articleId")

	// 校验数据
	if err != nil {
		beego.Error("获取文章id失败", err)
		this.Redirect("/showArticleList", 302)
		return
	}

	// 处理数据
	o := orm.NewOrm()
	var article models.Article
	article.Id = articleId
	err = o.Read(&article)
	if err != nil {
		beego.Error("获取数据失败", err)
		this.Redirect("/showArticleList", 302)
		return
	}

	// 返回视图
	this.Data["article"] = article
	this.Data["userName"] = this.GetSession("userName")
	this.TplName = "update.html"
}

// 处理添加文章
func (this *ArticleController) HandleUpdateArticle() {
	id, err := this.GetInt("articleId")
	// 获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	fileAddr := UploadFile(&this.Controller, "uploadname")

	// 校验数据
	if articleName == "" || content == "" || fileAddr == "" || err != nil {
		beego.Error("获取失败信息")
		this.Redirect("/showArticleList", 302)
		return
	}

	// 处理
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article)
	if err != nil {
		beego.Error("更新文章不存在", err)
		this.Redirect("/showArticleList", 302)
		return
	}

	article.Title = articleName
	article.Content = content
	article.Img = fileAddr

	_, err = o.Update(&article)
	if err != nil {
		beego.Error("更新文章失败", err)
		this.Redirect("/showArticleList", 302)
		return
	}

	// 返回
	this.Redirect("/showArticleList", 302)
}

// 封装文件上传检验函数
func UploadFile(this *beego.Controller, filePath string) string {
	// 图片上传
	file, head, err := this.GetFile(filePath)
	if err != nil {
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = "add.html"
		return ""
	}
	defer file.Close()

	// 检查文件大小
	if head.Size > 500000 {
		this.Data["errmsg"] = "文件太大"
		this.TplName = "add.html"
		return ""
	}

	// 检查文件类型
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		this.Data["errmsg"] = "文件类型错误， 请重新上传"
		this.TplName = "add.html"
		return ""
	}

	// 防止重名
	fileName := strconv.Itoa(int(time.Now().UnixNano()))

	//保存图片
	err = this.SaveToFile(filePath, "./static/img/"+fileName+ext)
	if err != nil {
		beego.Error("文件保存失败", err)
		return ""
	}

	return "./static/img/" + fileName + ext
}

// 删除文章
func (this *ArticleController) DelArticle() {
	// 获取数据
	id, err := this.GetInt("articleId")

	// 校验数据
	if err != nil {
		beego.Error("获取文章id失败", err)
		this.Redirect("/showArticleList", 302)
		return
	}

	// 处理数据
	o := orm.NewOrm()
	var article models.Article
	article.Id = id

	_, err = o.Delete(&article)
	if err != nil {
		beego.Error("删除失败", err)
		this.Redirect("/showArticleList", 302)
		return
	}

	// 返回视图
	this.Redirect("/showArticleList", 302)
}

// 展示添加文章分类页面
func (this *ArticleController) ShowAddArticleType() {
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)

	this.Data["articleTypes"] = articleTypes
	this.Data["userName"] = this.GetSession("userName")
	this.TplName = "addType.html"
}

// 添加文章请求
func (this *ArticleController) HandleAddArticleType() {
	// 获取数据
	typeName := this.GetString("typeName")

	// 校验数据
	if typeName == "" {
		beego.Error("获取数据失败")
		this.Redirect("/addArticleType", 302)
		return
	}

	// 处理数据
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName

	_, err := o.Insert(&articleType)
	if err != nil {
		beego.Error("插入数据失败")
		this.Redirect("/addArticleType", 302)
		return
	}

	// 返回视图
	this.Redirect("/addArticleType", 302)
}

// 删除分类
func (this *ArticleController) DelType() {
	// 获取数据
	typeId, err := this.GetInt("typeId")

	// 校验数据
	if err != nil {
		beego.Error("获取ID失败")
		this.Redirect("/addArticleType", 302)
		return
	}

	// 处理数据
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id = typeId

	_, err = o.Delete(&articleType)
	if err != nil {
		beego.Error("删除数据失败")
		this.Redirect("/addArticleType", 302)
		return
	}

	// 返回视图
	this.Redirect("/addArticleType", 302)
}

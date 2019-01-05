package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// 一对一		rel(one) reverse(one)
// 一对多		多的设置rel(fk)  一的设置reverse(many)
// 多对多		rel(m2m) reverse(many)
type User struct {
	Id       int
	UserName string `orm:"unique"`
	PassWord string

	Articles []*Article `orm:"reverse(many)"`
}

type Article struct {
	Id      int       `orm:"pk;auto"`
	Title   string    `orm:"size(100)"`
	Time    time.Time `orm:"type(datetime);auto_now"`
	Content string    `orm:"size(200)"`
	Count   int       `orm:"default(0)"`
	Img     string    `orm:"null"`

	ArticleType *ArticleType `orm:"rel(fk);null;on_delete(do_nothing)"`
	Users       []*User      `orm:"rel(m2m)"`
}

type ArticleType struct {
	Id       int
	TypeName string `orm:"size(20);unique"`

	Articles []*Article `orm:"reverse(many)"`
}

func init() {
	// 注册数据库
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/article?charset=utf8")
	// 注册表
	orm.RegisterModel(new(User), new(Article), new(ArticleType))

	// 运行
	orm.RunSyncdb("default", false, true)
}

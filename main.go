package main

import (
	_ "Article/models"
	_ "Article/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.AddFuncMap("pre", PrePage)
	beego.AddFuncMap("next", NextPage)
	beego.Run()
}

func PrePage(pageIndex int) int {
	if pageIndex == 1 {
		return 1
	}else {
		return pageIndex - 1
	}
}

func NextPage(pageIndex int, pageCount float64) int {
	if pageIndex == int(pageCount) {
		return int(pageCount)
	}else {
		return pageIndex + 1
	}

}
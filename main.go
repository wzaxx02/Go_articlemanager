package main

import (
	"web_blog/routes"

	"web_blog/model"
)

func main() {
	// 引用数据库
	model.InitDb()
	// 引入路由组件
	routes.InitRouter()

}

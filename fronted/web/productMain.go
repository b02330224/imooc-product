package main

import (
	"github.com/kataras/iris"
)

func main() {
	//创建iris实例
	app := iris.New()

	app.StaticWeb("/public", "./fronted/web/public")
	app.StaticWeb("/html", "./fronted/web/htmlProductShow")

	//启动服务
	app.Run(iris.Addr("0.0.0.0:80"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
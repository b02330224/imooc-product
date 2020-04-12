package main

import (
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"imooc-product/common"
	"imooc-product/fronted/middlerware"
	"imooc-product/fronted/web/controllers"
	"imooc-product/rabbitmq"
	"imooc-product/repositories"
	"imooc-product/services"
	"log"
)

func main() {
	//创建iris实例
	app := iris.New()

	//设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")

	//注册模版
	template := iris.HTML("./fronted/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)

	//设置模版目标
	app.StaticWeb("/public", "./fronted/web/public")
	app.StaticWeb("/html", "./fronted/web/htmlProductShow")

	//出现异常跳转导指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	//注册控制器
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Fatal(err)
	}


	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userRepository := repositories.NewUserRepository("user", db)
	userService := services.NewService(userRepository)
	user := mvc.New(app.Party("/user"))
	user.Register(ctx, userService)
	user.Handle(new(controllers.UserController))


	rabbitmq := rabbitmq.NewRabbitMQSimple("imoocProduct")

	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)

	order := repositories.NewOrderManagerRepository("`order`", db)
	orderService := services.NewOrderService(order)

	proProduct := app.Party("/product")
	proProduct.Use(middlerware.AuthConProduct)

	pro := mvc.New(proProduct)
	pro.Register(ctx, productService, orderService, rabbitmq)
	pro.Handle(new(controllers.ProductController))

	//启动服务
	app.Run(iris.Addr("0.0.0.0:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}

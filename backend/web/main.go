package main

import (
	gcontext "context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/mvc"
	"imooc-product/backend/web/controllers"
	"imooc-product/common"
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
	template := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)

	//设置模版目标
	app.StaticWeb("/assets", "./backend/web/assets")

	//出现异常跳转导指定页面
	app.OnAnyErrorCode(func(ctx context.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	//注册控制器
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Fatal(err)
	}


	ctx, cancel := gcontext.WithCancel(gcontext.Background())
	defer cancel()


	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderManagerRepository("`order`", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	//启动服务
	app.Run(iris.Addr("localhost:8081"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}

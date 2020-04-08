package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"html/template"
	"imooc-product/datamodels"
	"imooc-product/services"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type ProductController struct {
	Ctx iris.Context
	ProductService services.IProductService
	OrderService services.IOrderService
}

var (
	htmlOutPath = "./fronted/web/htmlProductShow/"
	templatePath = "./fronted/web/views/template/"
)

func (p *ProductController) GetGenerateHtml() {
	productIdString := p.Ctx.URLParam("productId")
	productId , err := strconv.Atoi(productIdString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}


	contentTmp, err := template.ParseFiles(filepath.Join(templatePath, "product.html"))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")
	product, err := p.ProductService.GetProductById(int64(productId))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	generateStaticHtml(p.Ctx, contentTmp, fileName, product)

}

func generateStaticHtml(ctx iris.Context, template *template.Template, fileName string, product *datamodels.Product)  {
	if exist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			ctx.Application().Logger().Error(err)
		}
	}

	file, err := os.OpenFile(fileName, os.O_CREATE | os.O_WRONLY, os.ModePerm)
	if err != nil {
		ctx.Application().Logger().Error(err)
	}
	defer file.Close()

	log.Println("product:", product)
	template.Execute(file, product)

}

func exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func (p *ProductController) GetDetail() mvc.View {
	product, err := p.ProductService.GetProductById(1)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Layout:"shared/productLayout.html",
		Name:"product/view.html",
		Data:iris.Map {
			"product" : product,
		},
	}
}

func (p *ProductController) GetOrder() mvc.View {
	productIdString := p.Ctx.URLParam("productId")
	userIdString := p.Ctx.GetCookie("uid")

	productId , err := strconv.Atoi(productIdString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	product, err := p.ProductService.GetProductById(int64(productId))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	var orderId int64
	showMessage := "抢购失败!"
	if product.ProductNum > 0 {
		//扣除商品数量
		product.ProductNum -= 1
		err := p.ProductService.UpdateProduct(product)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}

		//创建订单
		userId , err := strconv.Atoi(userIdString)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}

		order := &datamodels.Order{
			UserId:int64(userId),
			ProductId:int64(productId),
			OrderStatus:datamodels.OrderSuccess,
		}


		orderId , err = p.OrderService.InsertOrder(order)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		} else {
			showMessage = "抢购成功"
		}
	}

	return mvc.View{
		Layout:"shared/productLayout.html",
		Name:"product/result.html",
		Data:iris.Map {
			"orderId" : orderId,
			"showMessage" : showMessage,
		},
	}
}

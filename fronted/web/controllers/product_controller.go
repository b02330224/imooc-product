package controllers

import (
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"html/template"
	"imooc-product/datamodels"
	"imooc-product/rabbitmq"
	"imooc-product/services"
	"os"
	"path/filepath"
	"strconv"
)

type ProductController struct {
	Ctx iris.Context
	ProductService services.IProductService
	OrderService services.IOrderService
	RabbitMQ   *rabbitmq.RabbitMQ
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

func (p *ProductController) GetOrder() []byte {
	productIdString := p.Ctx.URLParam("productId")
	userIdString := p.Ctx.GetCookie("uid")
	productId, err := strconv.ParseInt(productIdString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	userId , err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	message := datamodels.NewMessage(userId, productId)
	byteMessage, err := json.Marshal(message)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	err = p.RabbitMQ.PublishSimple(string(byteMessage))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return []byte("true")
}

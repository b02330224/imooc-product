package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"imooc-product/datamodels"
	"imooc-product/services"
	"strconv"
)

type ProductController struct {
	Ctx iris.Context
	ProductService services.IProductService
	OrderService services.IOrderService
	Session *sessions.Session
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

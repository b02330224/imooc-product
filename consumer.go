package main

import (
	"fmt"
	"imooc-product/common"
	"imooc-product/rabbitmq"
	"imooc-product/repositories"
	"imooc-product/services"
)

func main() {
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}

	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)

	order := repositories.NewOrderManagerRepository("`order`", db)
	orderService := services.NewOrderService(order)

	rabbitConsumeSimple := rabbitmq.NewRabbitMQSimple("imoocProduct")

	rabbitConsumeSimple.ConsumeSimple(orderService, productService)
}

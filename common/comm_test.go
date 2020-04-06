package common

import (
	"fmt"
	"imooc-product/datamodels"
	"log"
	"testing"
)

func TestDataToStructByTagSql(t *testing.T) {
	data := map[string]string{
		"id" : "1",
		"productName" : "imooc 测试结构体",
		"productNum" : "2",
		"productImage" : "123",
		"productUrl" : "http://url",
	}
	log.Println(data)

	product := &datamodels.Product{}

	DataToStructByTagSql(data, product)

	fmt.Println(product)
}

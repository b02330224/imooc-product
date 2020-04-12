package services

import (
	"imooc-product/datamodels"
	"imooc-product/repositories"
)

type IOrderService interface {
	GetOrderById(int64) (*datamodels.Order, error)
	DeleteOrderById(int64) bool
	UpdateOrder(*datamodels.Order) error
	InsertOrder(*datamodels.Order) (int64, error)
	GetAllOrder() ([]*datamodels.Order, error)
	GetAllOrderInfo()(map[int]map[string]string, error)
	InsertOrderByMessage(message *datamodels.Message) (int64, error)
}

func NewOrderService(repository repositories.IOrderRepository) IOrderService {
	return &OrderService{orderRepository:repository}
}

type OrderService struct {
	orderRepository repositories.IOrderRepository
}

func (o *OrderService) GetOrderById(orderId int64) (order *datamodels.Order,err error) {
	return o.orderRepository.SelectByKey(orderId)
}

func (o *OrderService) DeleteOrderById(orderId int64) bool {
	isOk := o.orderRepository.Delete(orderId)
	return isOk
}

func (o *OrderService) UpdateOrder(order *datamodels.Order) error {
	return o.orderRepository.Update(order)
}

func (o *OrderService) InsertOrder(order *datamodels.Order) (int64, error) {
	return o.orderRepository.Insert(order)
}

func (o *OrderService) GetAllOrder() ([]*datamodels.Order, error) {
	return o.orderRepository.SelectAll()
}

func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
	return o.orderRepository.SelectAllWithInfo()
}

func (o *OrderService) InsertOrderByMessage(message *datamodels.Message) (orderId int64, err error) {
	order := &datamodels.Order{
		UserId:      message.UserId,
		ProductId:   message.ProductId,
		OrderStatus: datamodels.OrderSuccess,
	}
	return o.orderRepository.Insert(order)
}



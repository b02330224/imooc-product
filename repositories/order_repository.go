package repositories

import (
	"database/sql"
	"imooc-product/common"
	"imooc-product/datamodels"
	"strconv"
)

type IOrderRepository interface {
	Conn() error
	Insert(order *datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll()([]*datamodels.Order, error)
	SelectAllWithInfo()(map[int]map[string]string, error)
}

type OrderManagerRepository struct {
	table string
	mysqlConn *sql.DB
}

func (o *OrderManagerRepository) Conn() (err error) {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}

		o.mysqlConn = mysql
	}

	if o.table == "" {
		o.table = "order"
	}

	return
}

func (o *OrderManagerRepository) Insert(order *datamodels.Order) (orderId int64, err error) {
	if err = o.Conn();err != nil {
		return
	}

	sql := "insert into " + o.table + "(userId, productId, orderStatus) values (?, ?, ?)"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}

	result, err := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if err != nil {
		return
	}

	orderId, err = result.LastInsertId()
	return
}

func (o *OrderManagerRepository) Delete(orderId int64) (isOk bool) {
	if err := o.Conn();err != nil {
		return
	}
	sql := "delete from " + o.table + " where id = ?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}

	_, err = stmt.Exec(orderId)
	if err != nil {
		return false
	}

	return true
}

func (o *OrderManagerRepository) Update(order *datamodels.Order) error {
	if err := o.Conn();err != nil {
		return err
	}

	sql := "update " + o.table + " set userId=?, productId=?, orderStatus=? where id=" + strconv.FormatInt(order.Id, 10)
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderManagerRepository) SelectByKey(orderId int64) (order *datamodels.Order,err error) {
	if err := o.Conn();err != nil {
		return &datamodels.Order{}, err
	}

	sql := "select * from " + o.table + " where id" + strconv.FormatInt(orderId, 10)
	row, err := o.mysqlConn.Query(sql)
	if err != nil {
		return &datamodels.Order{}, err
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, err
	}

	common.DataToStructByTagSql(result, order)

	return
}

func (o *OrderManagerRepository) SelectAll() (orderArray []*datamodels.Order,err error) {
	if err := o.Conn();err != nil {
		return nil, err
	}

	sql := "select * from " + o.table
	rows, err := o.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}

	result := common.GetResultRows(rows)

	if len(result) == 0 {
		return nil, err
	}

	for _, v := range result {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return
}

func (o *OrderManagerRepository) SelectAllWithInfo() (orderMap map[int]map[string]string,err error) {
	if err := o.Conn();err != nil {
		return nil, err
	}

	sql := "select o.id, p.productName,o.orderStatus from `order` as o left join `product` as p on p.id = o.productId"
	rows, err := o.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}

	return common.GetResultRows(rows), err
}

func NewOrderManagerRepository(table string, sql *sql.DB) IOrderRepository {
	return &OrderManagerRepository{table:table, mysqlConn:sql}
}

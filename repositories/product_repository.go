package repositories

import (
	"database/sql"
	"imooc-product/common"
	"imooc-product/datamodels"
	"strconv"
)

//先开发接口
//实现定义的接口

type IProduct interface {
	//连接数据库
	Conn() error
	Insert(*datamodels.Product)(int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64)(*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}

type ProductManager struct {
	table string
	mysqlConn *sql.DB
}

func NewProductManager(table string, db *sql.DB) IProduct {
    return &ProductManager{table:table,mysqlConn:db}
}

func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}

	if p.table == "" {
		p.table = "product"
	}

	return
}

func (p *ProductManager) Insert(product *datamodels.Product) (productId int64,err error) {
	if err = p.Conn();err != nil {
		return
	}

	sql := "insert into product(productName, productNum, productImage, productUrl) values (?, ?, ?, ?)"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}

	result, err := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return
	}

	productId, err = result.LastInsertId()
	return
}

func (p *ProductManager) Delete(productId int64) bool {
	if err := p.Conn();err != nil {
		return false
	}

	sql := "delete from product where id = ?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}

	_, err = stmt.Exec(productId)
	if err != nil {
		return false
	}
	return true
}

func (p *ProductManager) Update(product *datamodels.Product) error {
	if err := p.Conn();err != nil {
		return err
	}

	sql := "update product set productName = ? , productNum = ? , productImage = ?, productUrl = ? where id=" + strconv.FormatInt(product.Id, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductManager) SelectByKey(productId int64) (product *datamodels.Product, err error) {
	if err = p.Conn();err != nil {
		return &datamodels.Product{}, err
	}

	sql := "select * from " + p.table + " where id = " + strconv.FormatInt(productId, 10)
	row, err := p.mysqlConn.Query(sql)
	if err != nil {
		return &datamodels.Product{}, err
	}
	defer row.Close()

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}

	product = &datamodels.Product{}
	common.DataToStructByTagSql(result, product)
	return product, nil
}

func (p *ProductManager) SelectAll() ([]*datamodels.Product,error) {
	if err := p.Conn();err != nil {
		return nil, err
	}

	sql := "select * from " + p.table
	rows,err := p.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, nil
	}

	productArray := make([]*datamodels.Product, 0)
	for _, v := range result {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)

		productArray = append(productArray, product)
	}

	return productArray, nil
}